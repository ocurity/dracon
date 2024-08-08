
.
$7b455932-99b1-4bca-8514-4dbf501767ba�ܰ�
pip-safety�
pyyaml:3.13Vulnerable Dependency
pyyaml[<4]0:�Advisory: Pyyaml before 4 uses ``yaml.load`` which has been assigned CVE-2017-18342.
Fixed Versions: [],Resources: [], More Info: https://data.safetycli.com/v/36333/97cBunknownJCVE-2017-18342RI7b455932-99b1-4bca-8514-4dbf501767ba:bbe502e1-5a8a-480d-9af4-43631777e759�
pyyaml:3.13Vulnerable Dependencypyyaml[<5.3.1]0:�Advisory: Pyyaml 5.3.1 includes a fix for CVE-2020-1747: A vulnerability was discovered in the PyYAML library in versions before 5.3.1, where it is susceptible to arbitrary code execution when it processes untrusted YAML files through the full_load method or with the FullLoader loader. Applications that use the library to process untrusted input may be vulnerable to this flaw. An attacker could use this flaw to execute arbitrary code on the system by abusing the python/object/new constructor.
Fixed Versions: [],Resources: [], More Info: https://data.safetycli.com/v/38100/97cBunknownJ
CVE-2020-1747RI7b455932-99b1-4bca-8514-4dbf501767ba:8dbcde8b-ebef-461d-8c46-5933f0c85274�
pyyaml:3.13Vulnerable Dependencypyyaml[<5.4]0:�Advisory: Pyyaml version 5.4 includes a fix for CVE-2020-14343: A vulnerability was discovered in the PyYAML library in versions before 5.4, where it is susceptible to arbitrary code execution when it processes untrusted YAML files through the full_load method or with the FullLoader loader. Applications that use the library to process untrusted input may be vulnerable to this flaw. This flaw allows an attacker to execute arbitrary code on the system by abusing the python/object/new constructor. This flaw is due to an incomplete fix for CVE-2020-1747.
https://bugzilla.redhat.com/show_bug.cgi?id=1860466
Fixed Versions: [],Resources: [], More Info: https://data.safetycli.com/v/39611/97cBunknownJCVE-2020-14343RI7b455932-99b1-4bca-8514-4dbf501767ba:deb975e2-8ad6-4066-bc62-a0604539aebc�
jinja2:2.10Vulnerable Dependencyjinja2[>=0,<2.10.1]0:�Advisory: Jinja2 2.10.1 adds 'SandboxedEnvironment' to handle 'str.format_map' in order to prevent code execution through untrusted format strings.
https://github.com/pallets/jinja/commit/a2a6c930bcca591a25d2b316fcfd2d6793897b26
Fixed Versions: [],Resources: [], More Info: https://data.safetycli.com/v/54679/97cBunknownJCVE-2019-10906RI7b455932-99b1-4bca-8514-4dbf501767ba:fa85d66b-6e1a-4f45-ada4-37919e3415ac�
jinja2:2.10Vulnerable Dependencyjinja2[>=0]0:�Advisory: In Jinja2, the from_string function is prone to Server Side Template Injection (SSTI) where it takes the "source" parameter as a template object, renders it, and then returns it. The attacker can exploit it with {{INJECTION COMMANDS}} in a URI. 
NOTE: The maintainer and multiple third parties believe that this vulnerability isn't valid because users shouldn't use untrusted templates without sandboxing.
Fixed Versions: [],Resources: [], More Info: https://data.safetycli.com/v/70612/97cBunknownJ
CVE-2019-8341RI7b455932-99b1-4bca-8514-4dbf501767ba:7f4ad151-c505-45af-9437-7a15a6fd7e69�
jinja2:2.10Vulnerable Dependencyjinja2[<3.1.3]0:�Advisory: Jinja2 before 3.1.3 is affected by a Cross-Site Scripting vulnerability. Special placeholders in the template allow writing code similar to Python syntax. It is possible to inject arbitrary HTML attributes into the rendered HTML template. The Jinja 'xmlattr' filter can be abused to inject arbitrary HTML attribute keys and values, bypassing the auto escaping mechanism and potentially leading to XSS. It may also be possible to bypass attribute validation checks if they are blacklist-based.
Fixed Versions: [],Resources: [], More Info: https://data.safetycli.com/v/64227/97cBunknownJCVE-2024-22195RI7b455932-99b1-4bca-8514-4dbf501767ba:1e8166ce-93bc-4d6b-8c0b-ab7939bed6f1�
jinja2:2.10Vulnerable Dependencyjinja2[<2.11.3]0:�Advisory: This affects the package jinja2 from 0.0.0 and before 2.11.3. The ReDoS vulnerability is mainly due to the '_punctuation_re regex' operator and its use of multiple wildcards. The last wildcard is the most exploitable as it searches for trailing punctuation. This issue can be mitigated by Markdown to format user content instead of the urlize filter, or by implementing request timeouts and limiting process memory.
Fixed Versions: [],Resources: [], More Info: https://data.safetycli.com/v/39525/97cBunknownJCVE-2020-28493RI7b455932-99b1-4bca-8514-4dbf501767ba:c73d960e-7b9f-426a-a9e4-56eac749d359�
idna:2.8Vulnerable Dependency
idna[<3.7]0:�Advisory: CVE-2024-3651 impacts the idna.encode() function, where a specially crafted argument could lead to significant resource consumption, causing a denial-of-service. In version 3.7, this function has been updated to reject such inputs efficiently, minimizing resource use. A practical workaround involves enforcing a maximum domain name length of 253 characters before encoding, as the vulnerability is triggered by unusually large inputs that normal operations wouldn't encounter.
Fixed Versions: [],Resources: [], More Info: https://data.safetycli.com/v/67895/97cBunknownJ
CVE-2024-3651RI7b455932-99b1-4bca-8514-4dbf501767ba:03b5ff48-5571-4eab-8cc1-60d5e6abeb6c�

hiredis:0.3.1Vulnerable Dependencyhiredis[<2.1.0]0:�Advisory: Hiredis (python wrapper for hiredis) 2.1.0 supports hiredis 1.1.0, that includes a security fix.
https://github.com/redis/hiredis-py/pull/135
Fixed Versions: [],Resources: [], More Info: https://data.safetycli.com/v/53276/97cBunknownJCVE-2021-32765RI7b455932-99b1-4bca-8514-4dbf501767ba:dced5cc4-2136-4839-8520-5f822eec5cbd�

aiohttp:3.5.3Vulnerable Dependencyaiohttp[<3.9.0]0:�Advisory: Aiohttp 3.9.0 includes a fix for CVE-2023-49081: Improper validation made it possible for an attacker to modify the HTTP request (e.g. to insert a new header) or create a new HTTP request if the attacker controls the HTTP version. The vulnerability only occurs if the attacker can control the HTTP version of the request.
https://github.com/aio-libs/aiohttp/security/advisories/GHSA-q3qx-c6g2-7pw2
Fixed Versions: [],Resources: [], More Info: https://data.safetycli.com/v/62582/97cBunknownJCVE-2023-49081RI7b455932-99b1-4bca-8514-4dbf501767ba:bef0fcad-7093-4288-80d8-665cda75f374�

aiohttp:3.5.3Vulnerable Dependencyaiohttp[<3.9.0]0:�Advisory: Affected versions of aiohttp are vulnerable to an Improper Validation vulnerability. It is possible for an attacker to modify the HTTP request (e.g. insert a new header) or even create a new HTTP request if the attacker controls the HTTP method. The vulnerability occurs only if the attacker can control the HTTP method (GET, POST etc.) of the request. If the attacker can control the HTTP version of the request it will be able to modify the request (request smuggling).
Fixed Versions: [],Resources: [], More Info: https://data.safetycli.com/v/62583/97cBunknownJCVE-2023-49082RI7b455932-99b1-4bca-8514-4dbf501767ba:1450429f-2e1b-4c8b-93d5-124b6f2ef529�

aiohttp:3.5.3Vulnerable Dependencyaiohttp[>1.0.5,<3.9.2]0:�Advisory: The vulnerability lies in the improper configuration of static resource resolution when aiohttp is used as a web server. It occurs when the follow_symlinks option is enabled without proper validation, leading to directory traversal vulnerabilities. Unauthorized access to arbitrary files on the system could potentially occur. The affected versions are >1.0.5, and the issue was patched in version 3.9.2. As a workaround, it is advised to disable the follow_symlinks option outside of a restricted local development environment, especially in a server accepting requests from remote users. Using a reverse proxy server to handle static resources is also recommended.
https://github.com/aio-libs/aiohttp/commit/1c335944d6a8b1298baf179b7c0b3069f10c514b
Fixed Versions: [],Resources: [], More Info: https://data.safetycli.com/v/64642/97cBunknownJCVE-2024-23334RI7b455932-99b1-4bca-8514-4dbf501767ba:7968ef82-cf5e-4d2e-9787-a89f866e5f24�

aiohttp:3.5.3Vulnerable Dependencyaiohttp[<3.9.4]0:�Advisory: aiohttp is an asynchronous HTTP client/server framework for asyncio and Python. A XSS vulnerability exists on index pages for static file handling. This vulnerability is fixed in 3.9.4. We have always recommended using a reverse proxy server (e.g. nginx) for serving static files. Users following the recommendation are unaffected. Other users can disable `show_index` if unable to upgrade. See CVE-2024-27306.
Fixed Versions: [],Resources: [], More Info: https://data.safetycli.com/v/70630/97cBunknownJCVE-2024-27306RI7b455932-99b1-4bca-8514-4dbf501767ba:78fbec5a-8968-4313-afac-ddcf4b235de3�

aiohttp:3.5.3Vulnerable Dependencyaiohttp[<3.8.6]0:�Advisory: Aiohttp 3.8.6 updates vendored copy of  'llhttp' to v9.1.3 to include a security fix.
https://github.com/aio-libs/aiohttp/security/advisories/GHSA-pjjw-qhg8-p2p9
Fixed Versions: [],Resources: [], More Info: https://data.safetycli.com/v/61657/97cBunknownRI7b455932-99b1-4bca-8514-4dbf501767ba:265ccdf9-1783-408c-8b95-1d69b03d9116�

aiohttp:3.5.3Vulnerable Dependencyaiohttp[<3.8.6]0:�Advisory: Aiohttp 3.8.6 includes a fix for CVE-2023-47627: The HTTP parser in AIOHTTP has numerous problems with header parsing, which could lead to request smuggling. This parser is only used when AIOHTTP_NO_EXTENSIONS is enabled (or not using a prebuilt wheel).
https://github.com/aio-libs/aiohttp/security/advisories/GHSA-gfw2-4jvh-wgfg
Fixed Versions: [],Resources: [], More Info: https://data.safetycli.com/v/62326/97cBunknownJCVE-2023-47627RI7b455932-99b1-4bca-8514-4dbf501767ba:1955efb5-e200-4924-8ba2-ac3688016b31�

aiohttp:3.5.3Vulnerable Dependencyaiohttp[<=3.8.4]0:�Advisory: Aiohttp 3.8.5 includes a fix for CVE-2023-37276: Sending a crafted HTTP request will cause the server to misinterpret one of the HTTP header values leading to HTTP request smuggling.
https://github.com/aio-libs/aiohttp/commit/9337fb3f2ab2b5f38d7e98a194bde6f7e3d16c40
https://github.com/aio-libs/aiohttp/security/advisories/GHSA-45c4-8wx5-qw6w
Fixed Versions: [],Resources: [], More Info: https://data.safetycli.com/v/59725/97cBunknownJCVE-2023-37276RI7b455932-99b1-4bca-8514-4dbf501767ba:d0b61b47-3bd3-453c-9210-55a48f22dce3�

aiohttp:3.5.3Vulnerable Dependencyaiohttp[<3.7.4]0:�Advisory: Aiohttp 3.7.4 includes a fix for CVE-2021-21330: In aiohttp before version 3.7.4 there is an open redirect vulnerability. A maliciously crafted link to an aiohttp-based web-server could redirect the browser to a different website. It is caused by a bug in the 'aiohttp.web_middlewares.normalize_path_middleware' middleware. A workaround can be to avoid using 'aiohttp.web_middlewares.normalize_path_middleware' in your applications.
https://github.com/aio-libs/aiohttp/security/advisories/GHSA-v6wp-4m6f-gcjg
Fixed Versions: [],Resources: [], More Info: https://data.safetycli.com/v/39659/97cBunknownJCVE-2021-21330RI7b455932-99b1-4bca-8514-4dbf501767ba:f31222f3-3484-479c-9926-6f65c17da4a6�

aiohttp:3.5.3Vulnerable Dependencyaiohttp[<3.9.1]0:�Advisory: The aiohttp versions minor than 3.9. has a vulnerability that affects the Python HTTP parser used in the aiohttp library. It allows for minor differences in allowable character sets, which could lead to robust frame boundary matching of proxies to protect against the injection of additional requests. The vulnerability also allows 
 exceptions during validation that aren't handled consistently with other malformed inputs.
Fixed Versions: [],Resources: [], More Info: https://data.safetycli.com/v/64644/97cBunknownJCVE-2024-23829RI7b455932-99b1-4bca-8514-4dbf501767ba:7e438725-6968-4630-94de-cd3314ff51f4�


aiohttp:3.5.3Vulnerable Dependencyaiohttp[<3.8.0]0:�	Advisory: Aiohttp 3.8.0 includes a fix for CVE-2023-47641: Affected versions of aiohttp have a security vulnerability regarding the inconsistent interpretation of the http protocol. HTTP/1.1 is a persistent protocol, if both Content-Length(CL) and Transfer-Encoding(TE) header values are present it can lead to incorrect interpretation of two entities that parse the HTTP and we can poison other sockets with this incorrect interpretation. A possible Proof-of-Concept (POC) would be a configuration with a reverse proxy(frontend) that accepts both CL and TE headers and aiohttp as backend. As aiohttp parses anything with chunked, we can pass a chunked123 as TE, the frontend entity will ignore this header and will parse Content-Length. The impact of this vulnerability is that it is possible to bypass any proxy rule, poisoning sockets to other users like passing Authentication Headers, also if it is present an Open Redirect an attacker could combine it to redirect random users to another website and log the request.
https://github.com/aio-libs/aiohttp/security/advisories/GHSA-xx9p-xxvh-7g8j
Fixed Versions: [],Resources: [], More Info: https://data.safetycli.com/v/62327/97cBunknownJCVE-2023-47641RI7b455932-99b1-4bca-8514-4dbf501767ba:a84ca030-3da6-475d-8077-adc03c5a4a92�

aiohttp:3.5.3Vulnerable Dependencyaiohttp[<3.8.0]0:�Advisory: Aiohttp 3.8.0 adds validation of HTTP header keys and values to prevent header injection.
https://github.com/aio-libs/aiohttp/issues/4818
Fixed Versions: [],Resources: [], More Info: https://data.safetycli.com/v/42692/97cBunknownRI7b455932-99b1-4bca-8514-4dbf501767ba:b175837a-1d67-45e8-b811-cd44acbc619d�
aiohttp-jinja2:1.1.0Vulnerable Dependencyaiohttp-jinja2[<1.1.1]0:�Advisory: Aiohttp-jinja2 1.1.1 updates minimal supported 'Jinja2' version to 2.10.1 to include security fixes.
Fixed Versions: [],Resources: [], More Info: https://data.safetycli.com/v/37095/97cBunknownJ
CVE-2014-1402RI7b455932-99b1-4bca-8514-4dbf501767ba:9a3ba115-0cc3-427b-b734-d1b3e17fc40c�
aiohttp-jinja2:1.1.0Vulnerable Dependencyaiohttp-jinja2[<1.1.1]0:�Advisory: Aiohttp-jinja2 1.1.1 updates minimal supported 'Jinja2' version to 2.10.1 to include security fixes.
Fixed Versions: [],Resources: [], More Info: https://data.safetycli.com/v/44431/97cBunknownJCVE-2016-10745RI7b455932-99b1-4bca-8514-4dbf501767ba:5547c6ed-d33b-4c08-87cc-cde8a002a47e�
aiohttp-jinja2:1.1.0Vulnerable Dependencyaiohttp-jinja2[<1.1.1]0:�Advisory: Aiohttp-jinja2 1.1.1 updates minimal supported 'Jinja2' version to 2.10.1 to include security fixes.
Fixed Versions: [],Resources: [], More Info: https://data.safetycli.com/v/44432/97cBunknownJCVE-2019-10906RI7b455932-99b1-4bca-8514-4dbf501767ba:37891a4f-4799-46e7-bc02-3bc095ebf2ec
