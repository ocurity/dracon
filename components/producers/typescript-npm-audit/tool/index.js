const childProc = require('child_process');
const filehound = require('filehound');
const fs        = require('fs');
const path      = require('path');
const readline  = require('readline');
const spawn     = require('await-spawn');
const util      = require('util');

const asyncMkdir     = util.promisify(fs.mkdir);
const asyncReadFile  = util.promisify(fs.readFile);
const asyncWriteFile = util.promisify(fs.writeFile);

if (process.argv.length != 4) {
    console.error('Usage: %s %s [SRC_DIR] [DEST_DIR]', process.argv[0], process.argv[1]);
    process.exit(1);
}

const srcDir  = process.argv[2];
const destDir = process.argv[3];

const metadataFiles = [
    'npm-shrinkwrap.json',
    'package.json',
    'package-lock.json',
    'yarn.lock'
];

const npmEnvironment = {
    'npm_config_audit':           'false',
    'npm_config_fund':            'false',
    'npm_config_loglevel':        'warn',
    'npm_config_update_notifier': 'false',
};

async function yarnAudit(packageDir) {
    // A bug in yarn currently makes `yarn audit --json` unusable
    // (https://github.com/yarnpkg/yarn/issues/7404), so we need to extract the
    // registry's audit report from `yarn audit --verbose` manually instead
    const yarn = childProc.spawn(
        'yarn',
        ['audit', '--verbose', '--groups', 'dependencies,optionalDependencies,peerDependencies'],
        {
            cwd:   packageDir,
            stdio: ['ignore', 'pipe', process.stderr]
        }
    );

    const yarnReader = require('readline').createInterface({
        input:    yarn.stdout,
        terminal: false
    });

    var reading = false;
    var report;

    for await (const line of yarnReader) {
        var match = line.match(/^verbose [\d.]+ Audit Response: (?<start>.*)/);
        if (match) {
            report  = match.groups.start;
            reading = true;
        } else if (reading) {
            report += line;

            if (line == '}') {
                break;
            }
        }
    }

    return report ? JSON.parse(report) : null;
}

async function npmAudit(packageDir) {
    // `npm audit` exits with a non-zero return code if vulnerabilities are
    // reported - although the minimum vulnerability severity required to
    // trigger this behaviour is configurable, it can't be disabled altogether,
    // preventing us from checking the return code to see whether the audit
    // succeeded
    const npm = childProc.spawn(
        'npm',
        ['audit', '--json', '--only=prod'],
        {
            cwd:   packageDir,
            env:   npmEnvironment,
            stdio: ['ignore', 'pipe', process.stderr]
        }
    );

    const npmReader = require('readline').createInterface({
        input:    npm.stdout,
        terminal: false
    });

    var report = '';

    for await (const line of npmReader) {
        report += line;
    }

    return report ? JSON.parse(report) : null;
}

async function npmInstall(packageDir) {
    try {
        var packageJSON = JSON.parse(
            await asyncReadFile(path.join(packageDir, 'package.json'), {
                encoding: 'utf8'
            })
        );
    } catch (err) {
        console.error(`${packageDir}: failed to read package.json: ${err}`);
        return 0;
    }

    var packageModified = false;
    var installed       = 0;

    // `npm audit` will fail if package.json defines any devDependencies that
    // aren't listed in package-lock.json, even if we tell it to only audit the
    // production dependencies - since we're not going to install the
    // devDependencies, remove that section from package.json altogether
    if ('devDependencies' in packageJSON) {
        delete packageJSON.devDependencies;
        packageModified = true;
    }

    for (const list of ['dependencies', 'optionalDependencies', 'peerDependencies']) {
        if (!(list in packageJSON)) {
            continue;
        }

        console.log(`${packageDir}: installing packages in category '${list}'`);

        for (const [name, version] of Object.entries(packageJSON[list])) {
            console.log(`${packageDir}: installing package: ${name}@${version}`);

            try {
                await spawn(
                    'npm',
                    [
                        'install', '--package-lock-only', '--only=prod', '--legacy-peer-deps',
                        name + '@' + version
                    ],
                    {
                        cwd:   packageDir,
                        env:   npmEnvironment,
                        stdio: ['ignore', 'ignore', process.stderr],
                    }
                );

                installed++;
            } catch (err) {
                // If the package installation failed, the package needs to be
                // removed from the dependencies list, otherwise `npm audit`
                // will later fail on the basis that not all dependencies in
                // package.json are installed
                console.error(`${packageDir}: failed to install package ${name}@${version}; removing from ${list}`);
                delete packageJSON[list][name];
                packageModified = true;
            }
        }
    }

    if (packageModified) {
        try {
            await asyncWriteFile(
                path.join(packageDir, 'package.json'),
                JSON.stringify(packageJSON, null, 2),
                { encoding: 'utf8' }
            );
        } catch (err) {
            console.error(`${packageDir}: failed to write package.json: ${err}`);
            return 0;
        }
    }

    return installed;
}

function auditSummary(report, summaryIntro) {
    const levels = ['critical', 'high', 'moderate', 'low', 'info'];

    const detail = levels.reduce((acc, level) => {
        acc.output.push(report.metadata.vulnerabilities[level] + ' ' + level);
        acc.total += report.metadata.vulnerabilities[level];

        return acc;
    }, { output: [], total: 0 });

    console.log(`${summaryIntro}: ${detail.total} (` + detail.output.join(', ') + ')');

    return detail.total;
}

(async function () {
    try {
        asyncMkdir(destDir, { recursive: true });
    } catch (err) {
        console.error(`Could not create output directory ${destDir}: ${err}`);
        process.exit(1);
    }

    const packageDirs = filehound.create()
        .paths(srcDir)
        .match(['package.json'])
        .discard('node_modules')
        .findSync()
        .map(file => {
            console.info(`Found package file: ${file}`);

            var relDir     = path.dirname(path.relative(srcDir, file));
            var relDestDir = path.join(destDir, relDir);

            // Mirror package metadata files into an equivalent directory structure
            // beneath DEST_DIR
            fs.mkdirSync(relDestDir, { recursive: true });
            metadataFiles.forEach(meta => {
                const metaSrc  = path.join(srcDir, relDir, meta);
                const metaDest = path.join(relDestDir, meta);

                if (fs.existsSync(metaSrc)) {
                    fs.copyFileSync(metaSrc, metaDest);
                }
            });

            return relDestDir;
        });

    console.log('Auditing found packages for vulnerable dependencies');

    for (const dir of packageDirs) {
        // If both npm and yarn lock files exist for this package, audit both,
        // since either could be used during deployment of the package
        var audited = false;

        if (fs.existsSync(path.join(dir, 'yarn.lock'))) {
            console.info(`${dir}: auditing yarn-based dependencies`);

            const report = await yarnAudit(dir);
            if (report) {
                const total = auditSummary(report, `${dir}: advisories for yarn-based dependencies`);
                if (total > 0) {
                    try {
                        await asyncWriteFile(
                            path.join(dir, 'package.yarn-audit'),
                            JSON.stringify(report),
                            { encoding: 'utf8' }
                        );
                    } catch (err) {
                        console.error(`${packageDir}: failed to write package.yarn-audit: ${err}`);
                    }
                }
            }

            audited = true;
        }

        if (
               fs.existsSync(path.join(dir, 'npm-shrinkwrap.json'))
            || fs.existsSync(path.join(dir, 'package-lock.json'))
        ) {
            console.info(`${dir}: auditing npm-based dependencies`);

            const report = await npmAudit(dir);
            if (report) {
                const total = auditSummary(report, `${dir}: advisories for npm-based dependencies`);
                if (total > 0) {
                    try {
                        await asyncWriteFile(
                            path.join(dir, 'package.npm-audit'),
                            JSON.stringify(report),
                            { encoding: 'utf8' }
                        );
                    } catch (err) {
                        console.error(`${packageDir}: failed to write package.npm-audit: ${err}`);
                    }
                }
            }

            audited = true;
        }

        // If we haven't audited the package at all yet, install dependency
        // metadata with npm using the information in package.json and audit
        // that - installation isn't guaranteed to succeed (e.g. when the
        // package is OS/architecture-specific and doesn't support the OS/
        // architecture we're currently running on), so install them
        // individually and we'll audit whatever we can
        if (!audited) {
            console.info(`${dir}: installing dependencies`);
            if (await npmInstall(dir) > 0) {
                console.info(`${dir}: auditing dependencies`);
                const report = await npmAudit(dir);
                if (report) {
                    const total = auditSummary(report, `${dir}: advisories for dependencies`);
                    if (total > 0) {
                        try {
                            await asyncWriteFile(
                                path.join(dir, 'package.npm-audit'),
                                JSON.stringify(report),
                                { encoding: 'utf8' }
                            );
                        } catch (err) {
                            console.error(`${packageDir}: failed to write package.npm-audit: ${err}`);
                        }
                    }
                }
            } else {
                console.info(`${dir}: no dependencies installed; skipping audit`);
            }
        }
    }
})();
