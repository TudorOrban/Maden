[Setup]
AppName=Maden
AppVersion=1.0
DefaultDirName={pf}\Maden
DefaultGroupName=Maden
OutputDir=userdocs:Inno Setup Examples Output
OutputBaseFilename=maden_setup
Compression=lzma
SolidCompression=yes

[Files]
Source: "D:\projects\programming\Golang\Maden\docker-compose.yml"; DestDir: "{app}";
Source: "D:\projects\programming\Golang\Maden\installer\start-maden.bat"; DestDir: "{app}";

[Tasks]
Name: "desktopicon"; Description: "Create a &desktop icon"; GroupDescription: "Additional icons:"

[Icons]
Name: "{group}\Maden"; Filename: "{app}\start-maden.bat"
Name: "{commondesktop}\Maden"; Filename: "{app}\start-maden.bat"; Tasks: desktopicon

[Run]
Filename: "{app}\installer-docker.bat"; Description: "Installing Docker"; Flags: runhidden runascurrentuser; StatusMsg: "Installing Docker..."
Filename: "{app}\start-maden.bat"; Description: "Starting Maden"; Flags: postinstall runascurrentuser; StatusMsg: "Starting Maden..."
