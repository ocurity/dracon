## OWASP Zaproxy producer

This producer runs owasp zaproxy in default mode using it's zap-full-scan.py script. 
This runs zap against the target specified in the parameter `producer-owasp-zaproxy-target`.
If you want zaproxy running in a more complex scenario (proxied end to end tests or using an openapi definitions file) you need to modify the task accordingly.