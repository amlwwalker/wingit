#!/bin/sh
test -f WingIt.dmg && rm WingIt.dmg
create-dmg \
--volname "WingIt" \
--volicon "_guiinterface/icons/wingit.icns" \
--background "DMG-background.png" \
--window-pos 200 120 \
--window-size 800 500 \
--icon-size 80 \
--icon "WingIt.app" 710 95 \
--hide-extension "WingIt.app" \
--app-drop-link 710 290 \
"WingIt.dmg" \
"_guiinterface/deploy/darwin/"