# Peaches' Log-Base Generator
[![Video Showcase](https://img.youtube.com/vi/t0vj5CILqtY/hqdefault.jpg)](https://youtu.be/t0vj5CILqtY)  
A simple-ish CLI program intended to make in-game HUD customizations via log-base techniques more accessible. The program doesn't necessarily replace handwritten log-base customizations for more complex or more demanding customizations. Currently the Generator only supports direct log-base generation, if you'd like to learn what that means or learn more about how log-base works, checkout my writeup on it in the [Xhud Wiki](https://github.com/PapaPeach/xhud/wiki/Log-File-Customizations) or JarateKing's dive into [Basefile Script Integration](https://github.com/JarateKing/TF2-Hud-Reference/blob/master/1-APPENDIX/BasefileScriptIntegration.md).
# Usage
1. Make a backup of your HUD prior to running the generator.
   The generator won't edit any files without prompting you first, and it will never write over any file. But I can't guarantee the program will work as intended 100% of the time, so don't bank on it.
2. Download the [latest release]() of the log-base generator.
3. Place the **log-base-generator.exe** in your HUD's root folder.
   (Where **info.vdf** is located).
4. Run the executable and follow the prompts.
5. Open **YourHud/log-base-copypasta.txt** and copy the generated button commands to  into their respective custom buttons in your HUD.
### Note
Currently there is no backwards navigation in the program, I recommend generating in smallish batches to avoid major set-backs if / when a typo is made.

# Reporting Issues
**With any issue, details on the file you were trying to access and previous selections that led you there are helpful in recreating and remedying the issue.**  
Opening a GitHub issue for the Log-Base Generator is the easiest way for me to keep track of issues to fix. Alternatively if there is an issue that needs addressing quickly, there will **soon** be a channel in my [Peaches' HUDs Discord](https://discord.gg/HyZRVtp) where you can report issues.

# Road Map
|                    Feature                    | Description                                                                                                                                                                                                                                             |   Status   |
| :-------------------------------------------: | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :--------: |
|                   Bug fixes                   | Battling spaghetti monster.                                                                                                                                                                                                                             |   Always   |
|       Default customization assignment        | Prompt user for what customization should be generated as the default on initial HUD user launch.                                                                                                                                                       |  Priority  |
|            Button console feedback            | Add prompt for feedback / debug text to be assigned to each customization when an individual option is selected on the HUD.                                                                                                                             |  Priority  |
|              Hard-reload option               | Add prompt or command to enabled hard reloads if a specific customization which requires it is selected.                                                                                                                                                |  Planned   |
|              Improved navigation              | Improved navigation such as directory / file autocompletion or fuzzy find. Backwards navigation in generation steps, canceling inputs, quit program command, etc.                                                                                       |  Planned   |
|           Animation customizations            | Add support for animation customizations and other supported non .res files. Requires investigation                                                                                                                                                     |  Planned   |
|          Automatic config detection           | Detect files previously generated via the Generator and prompt the user if they would like to use that, skipping several manual steps.<br>Potentially expanding to detect any existing file with code for log-base customization, even if user written. |  Planned   |
|         Nonvolatile selection storage         | Add some sort of loadable log or save file for applied settings so that unfinished generations can be resumed or updates can be accelerated.                                                                                                            |  Planned   |
| Support for exceptionally long customizations | There is a maximum length of code that an alias can contain in the Source engine, the program should eventually be able to detect this edge case and subdivide the alias appropriately to prevent issues.                                               |  Planned   |
|    Generator version and version detection    | If the generator updates to be non-backwards compatible or improves substantially. The program should mark generated files with a version number to allow for automated updating.                                                                       |   Likely   |
|                 Linux support                 | Support for HUDs used on Linux. Requires investigation into different handling of #base paths and slash preferences.                                                                                                                                    |   Likely   |
|             Customization presets             | Presets for things like sizing, positioning, visibility, etc. intended to speed up generation of common uses.                                                                                                                                           | Considered |
|            Customization converter            | Ability for program to detect and convert existing customization types to log-base.                                                                                                                                                                     | Considered |
|         Indirect base-log generation          | Add support for indirect log-base customizations to the generator.                                                                                                                                                                                      |  Unlikely  |

# Pull Requests
**Any pull request should have as narrow a scope as reasonably possible and have a single intention. In otherwards, fix one bug, typo, etc. per requests. Additionally, changes should be well documented in their methodology, effect, etc. if they would like to be considered.**  
Pull requests for feature additions are currently discouraged as the code is still young and more likely than not going to be rewritten actively by me. Once the code has matured and stabalized I will begin considering feature additions and update the README to reflect this.
