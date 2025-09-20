# Empire Command & Control
Empire is a powerful post-exploitation and adversary emulation framework designed to aid Red Teams and Penetration Testers.
Built with flexibility and modularity in mind, Empire enables security professionals to conduct sophisticated operations with ease.

The Empire server is written in Python 3, providing a robust and extensible backend for managing compromised systems.
Operators can interact with the server using Starkiller, a graphical user interface (GUI) that enhances usability and management.

## Key Features
- [x] **Server/Client Architecture** – Supports multiplayer operations with remote client access.
- [x] **Multi-Client Support** – Choose between a GUI (Starkiller) or command-line interface.
- [x] **Fully Encrypted Communications** – Ensures secure C2 channels
- [x] **Diverse Listener Support** – Communicate over HTTP/S, Malleable HTTP, and PHP.
- [x] **Extensive Module Library** – Over 400 tools in PowerShell, C#, and Python for post-exploitation and lateral movement.
- [x] **Donut Integration** – Generate shellcode for execution.
- [x] **Modular Plugin Interface** – Extend Empire with custom server features.
- [x] **Flexible Module Framework** – Easily add new capabilities.
- [x] **Advanced Obfuscation** – Integrated [ConfuserEx 2](https://github.com/mkaring/ConfuserEx) and [Invoke-Obfuscation](https://github.com/danielbohannon/Invoke-Obfuscation) for stealth.
- [x] **In-Memory Execution** – Load and execute .NET assemblies without touching disk.
- [x] **Customizable Bypasses** – Evade detection using JA3/S and JARM evasion techniques.
- [x] **MITRE ATT&CK Integration** – Map techniques and tactics directly to the framework.
- [x] **Built-in Roslyn Compiler** – Compile C# payloads on the fly (thanks to Covenant).
- [x] **Broad Deployment Support** – Install on Docker, Kali Linux, Ubuntu, and Debian.


Please see our Releases or Changelog page for detailed release notes.

Quickstart
When cloning this repository, you will need to recurse submodules.

git clone --recursive https://github.com/BC-SECURITY/Empire.git
Check out the Installation Page for install instructions.

Note: The main branch is a reflection of the latest changes and may not always be stable. After cloning the repo, you can checkout the latest stable release by running the setup/checkout-latest-tag.sh script.

git clone --recursive https://github.com/BC-SECURITY/Empire.git
cd Empire
./setup/checkout-latest-tag.sh
./ps-empire install -y
If you are using the sponsors version of Empire, it will pull the sponsors version of Starkiller. Because these are private repositories, you need to have ssh credentials configured for GitHub. Instructions can be found here.

Server
# Start Server
./ps-empire server

# Help
./ps-empire server -h
Check out the Empire Docs for more instructions on installing and using with Empire. For a complete list of changes, see the changelog.

Starkiller
