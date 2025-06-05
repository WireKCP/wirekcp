; NSIS Script to package a Windows executable and wintun.dll
; Enhanced with Modern UI (MUI2) for a better interface and an installer icon.

; --- General Installer Settings ---
Name "WireKCP Installer" ; Name of the installer displayed to the user
OutFile "WireKCP-Installer-Setup-x86_64.exe" ; Name of the generated installer executable

; Set the default installation directory.
; $PROGRAMFILES is a standard NSIS variable pointing to "C:\Program Files" or "C:\Program Files (x86)".
InstallDir "$PROGRAMFILES64\WireKCP"

; Request administrator privileges. This is crucial for writing to C:\Windows\System32.
RequestExecutionLevel admin

; Set the installer icon. Replace "my_app_icon.ico" with your actual icon file.
Icon "wirekcp.ico"
UninstallIcon "wirekcp.ico" ; Use the same icon for the uninstaller

; --- Modern UI (MUI2) Settings ---
!include "MUI2.nsh" ; Include the Modern UI header

; Define custom text for various pages
!define MUI_WELCOMEPAGE_TITLE "Welcome to the WireKCP Setup Wizard"
!define MUI_WELCOMEPAGE_TEXT "This wizard will guide you through the installation of WireKCP. \
  $\n$\nIt is recommended that you close all other applications before continuing. \
  $\n$\nClick Next to continue."

!define MUI_FINISHPAGE_TITLE "Completing the WireKCP Setup Wizard"
!define MUI_FINISHPAGE_TEXT "My Application has been installed on your computer. \
  $\n$\nClick Finish to exit this wizard."

!define MUI_ICON "wirekcp.ico"
; Optional: Define a header image for the installer pages.
; The image should be a bitmap (.bmp) file, typically around 150x57 pixels.
!define MUI_HEADERIMAGE
!define MUI_HEADERIMAGE_BITMAP "${NSISDIR}\Contrib\Graphics\Header\nsis3-metro-right.bmp" ; Replace with your actual header image file.
!define MUI_HEADERIMAGE_RIGHT

; Optional: Define a side image for the welcome/finish pages.
; The image should be a bitmap (.bmp) file, typically around 164x314 pixels.
!define MUI_WELCOMEFINISHPAGE_BITMAP "${NSISDIR}\Contrib\Graphics\Wizard\nsis3-metro.bmp" ; Uncomment and replace if you want a side image.

; Show a warning if the user tries to abort the installation.
!define MUI_ABORTWARNING

; --- User Interface Pages ---
; Welcome page
!insertmacro MUI_PAGE_WELCOME

; License page (optional - uncomment if you have a license agreement)
; Make sure you have a 'LICENSE.txt' file in the same directory as your .nsi script.
; LicenseData "LICENSE.txt"
; !insertmacro MUI_PAGE_LICENSE

; Directory page (allows user to choose installation path)
!insertmacro MUI_PAGE_DIRECTORY

; Installation files page (shows progress)
!insertmacro MUI_PAGE_INSTFILES

; Finish page
!insertmacro MUI_PAGE_FINISH

; --- Uninstaller User Interface Pages ---
; Uninstaller welcome page
!insertmacro MUI_UNPAGE_WELCOME

; Uninstaller confirmation page
!insertmacro MUI_UNPAGE_CONFIRM

; Uninstaller installation files page
!insertmacro MUI_UNPAGE_INSTFILES

; Uninstaller finish page
!insertmacro MUI_UNPAGE_FINISH

; --- Language Selection ---
; You can add multiple languages if needed. English is default.
!insertmacro MUI_LANGUAGE "English"

; --- Sections ---
; This section defines what files to install and where.
Section "WireKCP (required)"

  ; Set the output path for the main application files.
  ; $INSTDIR is the directory chosen by the user or the default InstallDir.
  SetOutPath "$INSTDIR"

  ; Copy your main executable to the installation directory.
  ; Replace "wirekcp.exe" with the actual name of your executable.
  File "dist\wirekcp_windows_amd64_v1\wirekcp.exe"

  ; Copy wintun.dll to the installation directory.
  File "wintun.dll"

 ; Set EnVar to operate on HKLM (System) environment variables
  EnVar::SetHKLM

  ; Create a new system environment variable WIREKCP_HOME and set it to $INSTDIR
  ; Note: EnVar::AddValue is typically for list-like variables. For a single value,
  ; WriteRegStr is direct, but EnVar::AddValue can also set it if it doesn't exist.
  ; We'll use WriteRegStr for direct setting and then EnVar::Update for broadcast.
  WriteRegStr HKLM "SYSTEM\CurrentControlSet\Control\Session Manager\Environment" "WIREKCP_HOME" "$INSTDIR"
  ; Broadcast message to notify other applications of the environment change for WIREKCP_HOME
  EnVar::Update "HKLM" "WIREKCP_HOME"

  ; Add %WIREKCP_HOME% to the system PATH environment variable.
  ; EnVar::AddValue will append the value if it's not already present.
  EnVar::AddValue "Path" "%WIREKCP_HOME%"
  ; Broadcast message to notify other applications of the environment change for Path
  EnVar::Update "HKLM" "Path"

  ; Write the uninstaller. This creates an executable that can remove the application.
  ; It's placed in the main installation directory.
  WriteUninstaller "$INSTDIR\uninstall.exe"

SectionEnd

; --- Uninstaller Section ---
; This section defines how to uninstall the application.
Section "Uninstall"

  ; Delete the main executable from the installation directory.
  Delete "$INSTDIR\wirekcp.exe"

  ; Delete wintun.dll from the system directory.
  Delete "$INSTDIR\wintun.dll"
  
  ; Remove %WIREKCP_HOME% from the system PATH environment variable.
  EnVar::DeleteValue "Path" "%WIREKCP_HOME%"
  ; Broadcast message to notify other applications of the environment change for Path
  EnVar::Update "HKLM" "Path"

  ; Delete the WIREKCP_HOME system environment variable.
  EnVar::Delete "WIREKCP_HOME"
  ; Broadcast message to notify other applications of the environment change for WIREKCP_HOME
  EnVar::Update "HKLM" "WIREKCP_HOME"
  
  ; Delete the uninstaller itself.
  Delete "$INSTDIR\uninstall.exe"

  ; Remove the main application directory if it's empty.
  RMDir "$INSTDIR"

SectionEnd
