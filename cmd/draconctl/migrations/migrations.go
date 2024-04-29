package migrations

import (
	"context"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
	"time"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"

	"github.com/ocurity/dracon/pkg/k8s"
	"github.com/ocurity/dracon/pkg/manifests"
)

var migrationsCmdConfig = struct {
	connStr         string
	migrationsTable string
	runAsK8sJob     bool
	timeout         int
}{}

var migrationsAsK8sJobConfig = struct {
	image         string
	kubeContext   string
	kubeConfig    string
	namespace     string
	leaseLockName string
	ssa           string
	dryRun        bool
	inCluster     bool
}{}

var migrationsCmd = &cobra.Command{
	Use:     "migrations",
	Long:    "A set of subcommands for managing database migrations",
	GroupID: "top-level",
}

type cmdEntrypoint func(*cobra.Command, []string) error

func init() {
	migrationsCmd.AddGroup(&cobra.Group{
		ID:    "Migrations",
		Title: "Migration Management:",
	})
	migrationsCmd.PersistentFlags().StringVar(&migrationsCmdConfig.connStr, "url", "postgres://postgres:postgres@localhost:5432", "Connection URL for the database.")
	migrationsCmd.PersistentFlags().StringVar(&migrationsCmdConfig.migrationsTable, "migrations-table", "oss-migrations", "Migrations table to inspect.")
	migrationsCmd.PersistentFlags().BoolVar(&migrationsCmdConfig.runAsK8sJob, "as-k8s-job", false, "Run command as a job in K8s")
	migrationsCmd.PersistentFlags().IntVarP(&migrationsCmdConfig.timeout, "timeout", "t", 10, "Timeout for command")
	migrationsCmd.PersistentFlags().BoolVar(&migrationsAsK8sJobConfig.dryRun, "dry-run", false, "Print the Job manifest to stdout instead of deploying it")
	migrationsCmd.PersistentFlags().BoolVar(&migrationsAsK8sJobConfig.inCluster, "in-cluster", false, "Binary is running inside a pod")
	migrationsCmd.PersistentFlags().StringVar(&migrationsAsK8sJobConfig.kubeContext, "kubecontext", "", "Use a specific kube context to execute opeations")
	migrationsCmd.PersistentFlags().StringVar(&migrationsAsK8sJobConfig.kubeConfig, "kubeconfig", "", "Path to kube config file")
	migrationsCmd.PersistentFlags().StringVar(&migrationsAsK8sJobConfig.leaseLockName, "lease-lock", "migration-job-lock", "Name for the lease lock configmap to use")
	migrationsCmd.PersistentFlags().StringVarP(&migrationsAsK8sJobConfig.namespace, "namespace", "n", "default", "Namespace where the migration job will be deployed")
	migrationsCmd.PersistentFlags().StringVarP(&migrationsAsK8sJobConfig.image, "image", "i", "", "Image to use containing draconctl binary to run command")
	migrationsCmd.PersistentFlags().StringVar(&migrationsAsK8sJobConfig.ssa, "ssa-name", "draconctl", "Name to use for server-side apply")
	// migrationsCmd.Flags().Bool("provide-password", false, "Provide the password via a console")
}

func RegisterMigrationsSubcommands(rootCmd *cobra.Command) {
	migrationsCmd.AddCommand(inspectSubCmd)
	migrationsCmd.AddCommand(applySubCmd)
	migrationsCmd.AddCommand(revertSubCmd)

	rootCmd.AddCommand(migrationsCmd)
}

// entrypointWrapper wraps the migration commands and checks the state of the flags and the
// environment and will make sure to run the command either locally, or as a Job in the K8s cluster
// If the binary is invoked inside a pod,
func entrypointWrapper(f cmdEntrypoint) cmdEntrypoint {
	return func(cmd *cobra.Command, args []string) error {
		// set timeout on the command
		ctx, cancel := context.WithTimeout(cmd.Context(), time.Duration(migrationsCmdConfig.timeout)*time.Second)
		cmd.SetContext(ctx)
		defer cancel()

		if migrationsCmdConfig.runAsK8sJob {
			return deployMigrationJob(cmd, migrationsAsK8sJobConfig.ssa)
		}

		if migrationsAsK8sJobConfig.dryRun {
			return errors.Errorf("you can't use the `--%s` flag without the `%s` flag", cmd.Flag("dry-run").Name, cmd.Flag("as-k8s-job").Name)
		} else if migrationsAsK8sJobConfig.kubeContext != "" {
			return errors.Errorf("you can't use the `--%s` flag without the `%s` flag", cmd.Flag("kube-context").Name, cmd.Flag("as-k8s-job").Name)
		} else if migrationsAsK8sJobConfig.kubeConfig != "" {
			return errors.Errorf("you can't use the `--%s` flag without the `%s` flag", cmd.Flag("kube-config").Name, cmd.Flag("as-k8s-job").Name)
		} else if migrationsAsK8sJobConfig.image != "" {
			return errors.Errorf("you can't use the `--%s` flag without the `%s` flag", cmd.Flag("image").Name, cmd.Flag("as-k8s-job").Name)
		}

		if migrationsAsK8sJobConfig.inCluster {
			// binary has been invoked inside a pod, we need to setup the client accordingly
			restCfg, err := rest.InClusterConfig()
			if err != nil {
				return errors.Errorf("could not initialise in-cluster K8s client config: %w", err)
			}
			leaderLockCtx, cancel := context.WithTimeout(cmd.Context(), time.Duration(migrationsCmdConfig.timeout)*time.Second)
			defer cancel()
			return grabLeaderLock(leaderLockCtx, f, cmd, args, restCfg)
		}
		return f(cmd, args)
	}
}

// grabLeaderLock will grab the leader lease and then execute the actual command
func grabLeaderLock(ctx context.Context, f cmdEntrypoint, cmd *cobra.Command, args []string, restCfg *rest.Config) (err error) {
	var client *clientset.Clientset
	client, err = clientset.NewForConfig(restCfg)
	if err != nil {
		return err
	}

	leaderElectionCtx, leaderElectionCancel := context.WithCancel(ctx)

	// listen for interrupts or the Linux SIGTERM signal and cancel
	// our context, which the leader election code will observe and
	// step down. channel will be closed
	setErr := sync.Once{}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalCh
		setErr.Do(func() {
			err = errors.New("received SIGTERM signal")
		})
		leaderElectionCancel()
	}()
	// cleanup signal handler
	defer signal.Stop(signalCh)

	leaseId := uuid.NewString()

	// start the leader election code loop
	leaderelection.RunOrDie(leaderElectionCtx, leaderelection.LeaderElectionConfig{
		Lock: &resourcelock.LeaseLock{
			LeaseMeta: metav1.ObjectMeta{
				Name:      migrationsAsK8sJobConfig.leaseLockName,
				Namespace: migrationsAsK8sJobConfig.namespace,
			},
			Client: client.CoordinationV1(),
			LockConfig: resourcelock.ResourceLockConfig{
				Identity: leaseId,
			},
		},
		// IMPORTANT: you MUST ensure that any code you have that
		// is protected by the lease must terminate **before**
		// you call cancel. Otherwise, you could have a background
		// loop still running and another process could
		// get elected before your background loop finished, violating
		// the stated goal of the lease.
		ReleaseOnCancel: true,
		LeaseDuration:   60 * time.Second,
		RenewDeadline:   15 * time.Second,
		RetryPeriod:     5 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				// defer wg.Done()
				cmd.Println("Starting controller loop...")
				cmd.SetContext(ctx)
				cmdErr := f(cmd, args)
				setErr.Do(func() {
					err = cmdErr
				})
			},
			OnStoppedLeading: func() {
				setErr.Do(func() {
					err = errors.Errorf("stopped being lease holder: %w", leaderElectionCtx.Err())
				})
				leaderElectionCancel()
			},
			OnNewLeader: func(identity string) {
				// we're notified when new leader elected
				if identity == leaseId {
					// I just got the lock
					return
				}
				cmd.Printf("new leader elected: %s", identity)
			},
		},
	})

	return err
}

// generateMigrationJob generates a Job manifest
func generateMigrationJob(cmdName string) *batchv1.Job {
	if migrationsAsK8sJobConfig.image == "" {
		migrationsAsK8sJobConfig.image = "europe-west1-docker.pkg.dev/oc-dracon-saas/demo/ocurity/dracon/draconctl:latest"
	}

	migrationJob := batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind: "Job",
			APIVersion: batchv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dracon-migrations",
			Namespace: migrationsAsK8sJobConfig.namespace,
			Labels: labels.Set{
				"v1.dracon.ocurity.com": "migrations",
				"generator":             "draconctl",
			},
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "dracon-migrations",
							Image: migrationsAsK8sJobConfig.image,
							Args:  getCleanedUpArgs(cmdName, os.Args[1:]),
							ImagePullPolicy: corev1.PullAlways,
						},
					},
					RestartPolicy:      corev1.RestartPolicyNever,
					ServiceAccountName: "dracon-migrations",
				},
			},
		},
	}

	return &migrationJob
}

func getCleanedUpArgs(cmdName string, args []string) []string {
	cleanedUpPodArgs := []string{}
	argsToRemove := map[string]bool{
		"--as-k8s-job":  false,
		"--kubecontext": true,
		"--kubeconfig":  true,
		"--timeout": 	 true,
		"--dry-run":     false,
		"--image":       true,
		"-i":            true,
	}
	var skip, skipNextArgToo bool
	for _, arg := range args {
		if skipNextArgToo {
			skipNextArgToo = false
			continue
		}
		if skipNextArgToo, skip = argsToRemove[arg]; !skip {
			cleanedUpPodArgs = append(cleanedUpPodArgs, arg)
		}
		if arg == cmdName {
			cleanedUpPodArgs = append(cleanedUpPodArgs, "--in-cluster")
		}
	}

	return cleanedUpPodArgs
}

func deployMigrationJob(cmd *cobra.Command, ssa string) error {
	migrationJob := generateMigrationJob(cmd.Name())
	if migrationsAsK8sJobConfig.dryRun {
		if err := manifests.BatchV1ObjEncoder.Encode(migrationJob, cmd.OutOrStdout()); err != nil {
			return errors.Errorf("could not marshal job manifest: %w", err)
		}
		return nil
	}

	if migrationsAsK8sJobConfig.kubeConfig == "" {
		migrationsAsK8sJobConfig.kubeConfig = path.Join(os.Getenv("HOME"), ".kube/config")
	}
	restCfg, err := clientcmd.BuildConfigFromFlags("", migrationsAsK8sJobConfig.kubeConfig)
	if err != nil {
		return errors.Errorf("%s: could not initialise K8s client config with: %w", migrationsAsK8sJobConfig.kubeConfig, err)
	}

	client, err := k8s.NewTypedClientForConfig(restCfg, ssa)
	if err != nil {
		return err
	}

	if err = client.Apply(cmd.Context(), migrationJob, migrationsAsK8sJobConfig.namespace, false); err != nil {
		return errors.Errorf("could not create migration job: %w", err)
	}

	ctx, cancel := context.WithTimeout(cmd.Context(), time.Duration(migrationsCmdConfig.timeout)*time.Second)
	defer cancel()
	return k8s.WatchJobPodLogs(ctx, client, migrationsAsK8sJobConfig.namespace, migrationJob.Name, 5, cmd.OutOrStdout())
}
