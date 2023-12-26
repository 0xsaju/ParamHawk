# ParamHawk

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)


ParamHawk is a tool designed to automate the discovery and extraction of parameters from target URLs. It removes duplicates, trims URLs, and organizes them for further analysis.
```
   _____                          _    _                _    
  |  __\\                        | |  | |              | |   
  | |__) |_ _ _ __ __ _ _ __ ___ | |__| | __ ___      _| | __
  |  ___/ _' | '__/ _' | '_ ' _ \|  __  |/ _' \ \ /\ / / |/ /
  | |  | (_| | | | (_| | | | | | | |  | | (_| |\ V  V /|   < 
  |_|   \__,_|_|  \__,_|_| |_| |_|_|  |_|\__,_| \_/\_/ |_|\_\
                                                            
#################################################################################################################
# Tools Name: ParamHawk                                                                                         #
# Description: This script is designed to automate the discovery and extraction of parameters from target URLs. #
#              It removes duplicates, trims URLs, and organizes them for further analysis.                      #
# Author: 0xsaju                                                                                                #
# LinkedIn: https://linkedin.com/in/0xsaju                                                                      #
# Version: v_1.0                                                                                                #
#################################################################################################################

-c	Collect all URLs from the target using waybackurls
-d  Target domain
-p	Generate param URLs

```
## Features

- **Collect URLs:** Use waybackurls to collect all URLs from the target domain.
- **Generate Param URLs:** Extract and organize unique parameter URLs.

## Installation

```
go install github.com/0xsaju/paramhawk@latest
```

Make sure your Go bin directory is in your PATH.

```
paramhawk -d example.com -c -p
```

    -d: Target domain.
    -c: Collect all URLs from the target using waybackurls.
    -p: Generate param URLs.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
Author

    0xsaju - LinkedIn

## Acknowledgments

    Inspired by the need for automated parameter discovery in web security testing.
    Hat tip to the open-source community.


## Please note:

- Update the installation command with your GitHub username and repository name.
- Ensure the `LICENSE` file is present with the MIT License text.
- Personalize the "Author" section with your information.
- Consider adding more sections like "Contributing," "Issues," or "Support."

Feel free to customize the README further based on your preferenc
