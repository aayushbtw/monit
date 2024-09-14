Output monit.gif

Require echo

Set FontSize 18
Set Width 1920
Set Height 1080
Set Padding 30

Set Margin 80
Set MarginFill "#674EFF"
Set WindowBar Colorful
Set WindowBarSize 40
Set BorderRadius 8

Set TypingSpeed 100ms

Type "ssh [ip] -p 23234 -i [ssh_key]" Sleep 200ms
Hide
Backspace 30
Type "ssh localhost -p 23234 -i .ssh/host_monit_ed25519" Enter
Sleep 500ms
Show

Sleep 5s
Screenshot screenshot.png
