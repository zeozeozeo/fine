package fine

import (
	"github.com/veandco/go-sdl2/sdl"
)

func (app *App) handleKeyboardEvent(event *sdl.KeyboardEvent) {
	// keycode := string(event.Keysym.Sym)
	key := Key(event.Keysym.Sym)
	isKeyDown := app.IsKeyDown(key)

	if app.OnKeyEvent != nil {
		app.OnKeyEvent(key, KeyDirection(event.State), app)
	}

	if event.State == sdl.RELEASED {
		app.JustUpKeys = append(app.JustUpKeys, key)
	}

	if event.State == sdl.PRESSED && !isKeyDown {
		app.HeldKeys = append(app.HeldKeys, key)
		app.JustDownKeys = append(app.JustDownKeys, key)
	} else if event.State == sdl.RELEASED && isKeyDown {
		// Find key and remove it
		for idx, heldKey := range app.HeldKeys {
			if key == heldKey {
				app.HeldKeys = append(app.HeldKeys[:idx], app.HeldKeys[idx+1:]...)
				break
			}
		}
	}
}

// Checks if a key is currently pressed.
func (app *App) IsKeyDown(key Key) bool {
	for _, heldKey := range app.HeldKeys {
		if key == heldKey {
			return true
		}
	}
	return false
}

// Checks if all passed keys are currently pressed.
func (app *App) AreKeysDown(keys ...Key) bool {
	for _, key := range keys {
		if !app.IsKeyDown(key) {
			return false
		}
	}
	return true
}

// Checks if a key was pressed on this frame.
func (app *App) IsKeyJustDown(key Key) bool {
	for _, downKey := range app.JustDownKeys {
		if key == downKey { // donkey
			return true
		}
	}
	return false
}

// Checks if a key was released on this frame.
func (app *App) IsKeyJustUp(key Key) bool {
	for _, upKey := range app.JustUpKeys {
		if key == upKey {
			return true
		}
	}
	return false
}

// Keyboard keys.
type Key int          // Key.
type KeyDirection int // Key direction (up or down).

const (
	KEYDIR_UP   KeyDirection = 0 // Key is released.
	KEYDIR_DOWN KeyDirection = 1 // Key is pressed.
)

const (
	KEY_UNKNOWN Key = sdl.K_UNKNOWN // "" (no name, empty string)

	KEY_RETURN     Key = sdl.K_RETURN     // "Return" (the Enter key (main keyboard))
	KEY_ESCAPE     Key = sdl.K_ESCAPE     // "Escape" (the Esc key)
	KEY_BACKSPACE  Key = sdl.K_BACKSPACE  // "Backspace"
	KEY_TAB        Key = sdl.K_TAB        // "Tab" (the Tab key)
	KEY_SPACE      Key = sdl.K_SPACE      // "Space" (the Space Bar key(s))
	KEY_EXCLAIM    Key = sdl.K_EXCLAIM    // "!"
	KEY_QUOTEDBL   Key = sdl.K_QUOTEDBL   // """
	KEY_HASH       Key = sdl.K_HASH       // "#"
	KEY_PERCENT    Key = sdl.K_PERCENT    // "%"
	KEY_DOLLAR     Key = sdl.K_DOLLAR     // "$"
	KEY_AMPERSAND  Key = sdl.K_AMPERSAND  // "&"
	KEY_QUOTE      Key = sdl.K_QUOTE      // "'"
	KEY_LEFTPAREN  Key = sdl.K_LEFTPAREN  // "("
	KEY_RIGHTPAREN Key = sdl.K_RIGHTPAREN // ")"
	KEY_ASTERISK   Key = sdl.K_ASTERISK   // "*"
	KEY_PLUS       Key = sdl.K_PLUS       // "+"
	KEY_COMMA      Key = sdl.K_COMMA      // ","
	KEY_MINUS      Key = sdl.K_MINUS      // "-"
	KEY_PERIOD     Key = sdl.K_PERIOD     // "."
	KEY_SLASH      Key = sdl.K_SLASH      // "/"
	KEY_0          Key = sdl.K_0          // "0"
	KEY_1          Key = sdl.K_1          // "1"
	KEY_2          Key = sdl.K_2          // "2"
	KEY_3          Key = sdl.K_3          // "3"
	KEY_4          Key = sdl.K_4          // "4"
	KEY_5          Key = sdl.K_5          // "5"
	KEY_6          Key = sdl.K_6          // "6"
	KEY_7          Key = sdl.K_7          // "7"
	KEY_8          Key = sdl.K_8          // "8"
	KEY_9          Key = sdl.K_9          // "9"
	KEY_COLON      Key = sdl.K_COLON      // ":"
	KEY_SEMICOLON  Key = sdl.K_SEMICOLON  // ";"
	KEY_LESS       Key = sdl.K_LESS       // "<"
	KEY_EQUALS     Key = sdl.K_EQUALS     // "="
	KEY_GREATER    Key = sdl.K_GREATER    // ">"
	KEY_QUESTION   Key = sdl.K_QUESTION   // "?"
	KEY_AT         Key = sdl.K_AT         // "@"

	KEY_LEFTBRACKET  Key = sdl.K_LEFTBRACKET  // "["
	KEY_BACKSLASH    Key = sdl.K_BACKSLASH    // "\"
	KEY_RIGHTBRACKET Key = sdl.K_RIGHTBRACKET // "]"
	KEY_CARET        Key = sdl.K_CARET        // "^"
	KEY_UNDERSCORE   Key = sdl.K_UNDERSCORE   // "_"
	KEY_BACKQUOTE    Key = sdl.K_BACKQUOTE    // "`"
	KEY_a            Key = sdl.K_a            // "A"
	KEY_b            Key = sdl.K_b            // "B"
	KEY_c            Key = sdl.K_c            // "C"
	KEY_d            Key = sdl.K_d            // "D"
	KEY_e            Key = sdl.K_e            // "E"
	KEY_f            Key = sdl.K_f            // "F"
	KEY_g            Key = sdl.K_g            // "G"
	KEY_h            Key = sdl.K_h            // "H"
	KEY_i            Key = sdl.K_i            // "I"
	KEY_j            Key = sdl.K_j            // "J"
	KEY_k            Key = sdl.K_k            // "K"
	KEY_l            Key = sdl.K_l            // "L"
	KEY_m            Key = sdl.K_m            // "M"
	KEY_n            Key = sdl.K_n            // "N"
	KEY_o            Key = sdl.K_o            // "O"
	KEY_p            Key = sdl.K_p            // "P"
	KEY_q            Key = sdl.K_q            // "Q"
	KEY_r            Key = sdl.K_r            // "R"
	KEY_s            Key = sdl.K_s            // "S"
	KEY_t            Key = sdl.K_t            // "T"
	KEY_u            Key = sdl.K_u            // "U"
	KEY_v            Key = sdl.K_v            // "V"
	KEY_w            Key = sdl.K_w            // "W"
	KEY_x            Key = sdl.K_x            // "X"
	KEY_y            Key = sdl.K_y            // "Y"
	KEY_z            Key = sdl.K_z            // "Z"

	KEY_CAPSLOCK Key = sdl.K_CAPSLOCK // "CapsLock"

	KEY_F1  Key = sdl.K_F1  // "F1"
	KEY_F2  Key = sdl.K_F2  // "F2"
	KEY_F3  Key = sdl.K_F3  // "F3"
	KEY_F4  Key = sdl.K_F4  // "F4"
	KEY_F5  Key = sdl.K_F5  // "F5"
	KEY_F6  Key = sdl.K_F6  // "F6"
	KEY_F7  Key = sdl.K_F7  // "F7"
	KEY_F8  Key = sdl.K_F8  // "F8"
	KEY_F9  Key = sdl.K_F9  // "F9"
	KEY_F10 Key = sdl.K_F10 // "F10"
	KEY_F11 Key = sdl.K_F11 // "F11"
	KEY_F12 Key = sdl.K_F12 // "F12"

	KEY_PRINTSCREEN Key = sdl.K_PRINTSCREEN // "PrintScreen"
	KEY_SCROLLLOCK  Key = sdl.K_SCROLLLOCK  // "ScrollLock"
	KEY_PAUSE       Key = sdl.K_PAUSE       // "Pause" (the Pause / Break key)
	KEY_INSERT      Key = sdl.K_INSERT      // "Insert" (insert on PC, help on some Mac keyboards (but does send code 73, not 117))
	KEY_HOME        Key = sdl.K_HOME        // "Home"
	KEY_PAGEUP      Key = sdl.K_PAGEUP      // "PageUp"
	KEY_DELETE      Key = sdl.K_DELETE      // "Delete"
	KEY_END         Key = sdl.K_END         // "End"
	KEY_PAGEDOWN    Key = sdl.K_PAGEDOWN    // "PageDown"
	KEY_RIGHT       Key = sdl.K_RIGHT       // "Right" (the Right arrow key (navigation keypad))
	KEY_LEFT        Key = sdl.K_LEFT        // "Left" (the Left arrow key (navigation keypad))
	KEY_DOWN        Key = sdl.K_DOWN        // "Down" (the Down arrow key (navigation keypad))
	KEY_UP          Key = sdl.K_UP          // "Up" (the Up arrow key (navigation keypad))

	KEY_NUMLOCKCLEAR Key = sdl.K_NUMLOCKCLEAR // "Numlock" (the Num Lock key (PC) / the Clear key (Mac))
	KEY_KP_DIVIDE    Key = sdl.K_KP_DIVIDE    // "Keypad /" (the / key (numeric keypad))
	KEY_KP_MULTIPLY  Key = sdl.K_KP_MULTIPLY  // "Keypad *" (the * key (numeric keypad))
	KEY_KP_MINUS     Key = sdl.K_KP_MINUS     // "Keypad -" (the - key (numeric keypad))
	KEY_KP_PLUS      Key = sdl.K_KP_PLUS      // "Keypad +" (the + key (numeric keypad))
	KEY_KP_ENTER     Key = sdl.K_KP_ENTER     // "Keypad Enter" (the Enter key (numeric keypad))
	KEY_KP_1         Key = sdl.K_KP_1         // "Keypad 1" (the 1 key (numeric keypad))
	KEY_KP_2         Key = sdl.K_KP_2         // "Keypad 2" (the 2 key (numeric keypad))
	KEY_KP_3         Key = sdl.K_KP_3         // "Keypad 3" (the 3 key (numeric keypad))
	KEY_KP_4         Key = sdl.K_KP_4         // "Keypad 4" (the 4 key (numeric keypad))
	KEY_KP_5         Key = sdl.K_KP_5         // "Keypad 5" (the 5 key (numeric keypad))
	KEY_KP_6         Key = sdl.K_KP_6         // "Keypad 6" (the 6 key (numeric keypad))
	KEY_KP_7         Key = sdl.K_KP_7         // "Keypad 7" (the 7 key (numeric keypad))
	KEY_KP_8         Key = sdl.K_KP_8         // "Keypad 8" (the 8 key (numeric keypad))
	KEY_KP_9         Key = sdl.K_KP_9         // "Keypad 9" (the 9 key (numeric keypad))
	KEY_KP_0         Key = sdl.K_KP_0         // "Keypad 0" (the 0 key (numeric keypad))
	KEY_KP_PERIOD    Key = sdl.K_KP_PERIOD    // "Keypad ." (the . key (numeric keypad))

	KEY_APPLICATION    Key = sdl.K_APPLICATION    // "Application" (the Application / Compose / Context Menu (Windows) key)
	KEY_POWER          Key = sdl.K_POWER          // "Power" (The USB document says this is a status flag, not a physical key - but some Mac keyboards do have a power key.)
	KEY_KP_EQUALS      Key = sdl.K_KP_EQUALS      // "Keypad =" (the = key (numeric keypad))
	KEY_F13            Key = sdl.K_F13            // "F13"
	KEY_F14            Key = sdl.K_F14            // "F14"
	KEY_F15            Key = sdl.K_F15            // "F15"
	KEY_F16            Key = sdl.K_F16            // "F16"
	KEY_F17            Key = sdl.K_F17            // "F17"
	KEY_F18            Key = sdl.K_F18            // "F18"
	KEY_F19            Key = sdl.K_F19            // "F19"
	KEY_F20            Key = sdl.K_F20            // "F20"
	KEY_F21            Key = sdl.K_F21            // "F21"
	KEY_F22            Key = sdl.K_F22            // "F22"
	KEY_F23            Key = sdl.K_F23            // "F23"
	KEY_F24            Key = sdl.K_F24            // "F24"
	KEY_EXECUTE        Key = sdl.K_EXECUTE        // "Execute"
	KEY_HELP           Key = sdl.K_HELP           // "Help"
	KEY_MENU           Key = sdl.K_MENU           // "Menu"
	KEY_SELECT         Key = sdl.K_SELECT         // "Select"
	KEY_STOP           Key = sdl.K_STOP           // "Stop"
	KEY_AGAIN          Key = sdl.K_AGAIN          // "Again" (the Again key (Redo))
	KEY_UNDO           Key = sdl.K_UNDO           // "Undo"
	KEY_CUT            Key = sdl.K_CUT            // "Cut"
	KEY_COPY           Key = sdl.K_COPY           // "Copy"
	KEY_PASTE          Key = sdl.K_PASTE          // "Paste"
	KEY_FIND           Key = sdl.K_FIND           // "Find"
	KEY_MUTE           Key = sdl.K_MUTE           // "Mute"
	KEY_VOLUMEUP       Key = sdl.K_VOLUMEUP       // "VolumeUp"
	KEY_VOLUMEDOWN     Key = sdl.K_VOLUMEDOWN     // "VolumeDown"
	KEY_KP_COMMA       Key = sdl.K_KP_COMMA       // "Keypad ," (the Comma key (numeric keypad))
	KEY_KP_EQUALSAS400 Key = sdl.K_KP_EQUALSAS400 // "Keypad = (AS400)" (the Equals AS400 key (numeric keypad))

	KEY_ALTERASE   Key = sdl.K_ALTERASE   // "AltErase" (Erase-Eaze)
	KEY_SYSREQ     Key = sdl.K_SYSREQ     // "SysReq" (the SysReq key)
	KEY_CANCEL     Key = sdl.K_CANCEL     // "Cancel"
	KEY_CLEAR      Key = sdl.K_CLEAR      // "Clear"
	KEY_PRIOR      Key = sdl.K_PRIOR      // "Prior"
	KEY_RETURN2    Key = sdl.K_RETURN2    // "Return"
	KEY_SEPARATOR  Key = sdl.K_SEPARATOR  // "Separator"
	KEY_OUT        Key = sdl.K_OUT        // "Out"
	KEY_OPER       Key = sdl.K_OPER       // "Oper"
	KEY_CLEARAGAIN Key = sdl.K_CLEARAGAIN // "Clear / Again"
	KEY_CRSEL      Key = sdl.K_CRSEL      // "CrSel"
	KEY_EXSEL      Key = sdl.K_EXSEL      // "ExSel"

	KEY_KP_00              Key = sdl.K_KP_00              // "Keypad 00" (the 00 key (numeric keypad))
	KEY_KP_000             Key = sdl.K_KP_000             // "Keypad 000" (the 000 key (numeric keypad))
	KEY_THOUSANDSSEPARATOR Key = sdl.K_THOUSANDSSEPARATOR // "ThousandsSeparator" (the Thousands Separator key)
	KEY_DECIMALSEPARATOR   Key = sdl.K_DECIMALSEPARATOR   // "DecimalSeparator" (the Decimal Separator key)
	KEY_CURRENCYUNIT       Key = sdl.K_CURRENCYUNIT       // "CurrencyUnit" (the Currency Unit key)
	KEY_CURRENCYSUBUNIT    Key = sdl.K_CURRENCYSUBUNIT    // "CurrencySubUnit" (the Currency Subunit key)
	KEY_KP_LEFTPAREN       Key = sdl.K_KP_LEFTPAREN       // "Keypad (" (the Left Parenthesis key (numeric keypad))
	KEY_KP_RIGHTPAREN      Key = sdl.K_KP_RIGHTPAREN      // "Keypad )" (the Right Parenthesis key (numeric keypad))
	KEY_KP_LEFTBRACE       Key = sdl.K_KP_LEFTBRACE       // "Keypad {" (the Left Brace key (numeric keypad))
	KEY_KP_RIGHTBRACE      Key = sdl.K_KP_RIGHTBRACE      // "Keypad }" (the Right Brace key (numeric keypad))
	KEY_KP_TAB             Key = sdl.K_KP_TAB             // "Keypad Tab" (the Tab key (numeric keypad))
	KEY_KP_BACKSPACE       Key = sdl.K_KP_BACKSPACE       // "Keypad Backspace" (the Backspace key (numeric keypad))
	KEY_KP_A               Key = sdl.K_KP_A               // "Keypad A" (the A key (numeric keypad))
	KEY_KP_B               Key = sdl.K_KP_B               // "Keypad B" (the B key (numeric keypad))
	KEY_KP_C               Key = sdl.K_KP_C               // "Keypad C" (the C key (numeric keypad))
	KEY_KP_D               Key = sdl.K_KP_D               // "Keypad D" (the D key (numeric keypad))
	KEY_KP_E               Key = sdl.K_KP_E               // "Keypad E" (the E key (numeric keypad))
	KEY_KP_F               Key = sdl.K_KP_F               // "Keypad F" (the F key (numeric keypad))
	KEY_KP_XOR             Key = sdl.K_KP_XOR             // "Keypad XOR" (the XOR key (numeric keypad))
	KEY_KP_POWER           Key = sdl.K_KP_POWER           // "Keypad ^" (the Power key (numeric keypad))
	KEY_KP_PERCENT         Key = sdl.K_KP_PERCENT         // "Keypad %" (the Percent key (numeric keypad))
	KEY_KP_LESS            Key = sdl.K_KP_LESS            // "Keypad <" (the Less key (numeric keypad))
	KEY_KP_GREATER         Key = sdl.K_KP_GREATER         // "Keypad >" (the Greater key (numeric keypad))
	KEY_KP_AMPERSAND       Key = sdl.K_KP_AMPERSAND       // "Keypad &" (the & key (numeric keypad))
	KEY_KP_DBLAMPERSAND    Key = sdl.K_KP_DBLAMPERSAND    // "Keypad &&" (the && key (numeric keypad))
	KEY_KP_VERTICALBAR     Key = sdl.K_KP_VERTICALBAR     // "Keypad |" (the | key (numeric keypad))
	KEY_KP_DBLVERTICALBAR  Key = sdl.K_KP_DBLVERTICALBAR  // "Keypad ||" (the || key (numeric keypad))
	KEY_KP_COLON           Key = sdl.K_KP_COLON           // "Keypad :" (the : key (numeric keypad))
	KEY_KP_HASH            Key = sdl.K_KP_HASH            // "Keypad #" (the # key (numeric keypad))
	KEY_KP_SPACE           Key = sdl.K_KP_SPACE           // "Keypad Space" (the Space key (numeric keypad))
	KEY_KP_AT              Key = sdl.K_KP_AT              // "Keypad @" (the @ key (numeric keypad))
	KEY_KP_EXCLAM          Key = sdl.K_KP_EXCLAM          // "Keypad !" (the ! key (numeric keypad))
	KEY_KP_MEMSTORE        Key = sdl.K_KP_MEMSTORE        // "Keypad MemStore" (the Mem Store key (numeric keypad))
	KEY_KP_MEMRECALL       Key = sdl.K_KP_MEMRECALL       // "Keypad MemRecall" (the Mem Recall key (numeric keypad))
	KEY_KP_MEMCLEAR        Key = sdl.K_KP_MEMCLEAR        // "Keypad MemClear" (the Mem Clear key (numeric keypad))
	KEY_KP_MEMADD          Key = sdl.K_KP_MEMADD          // "Keypad MemAdd" (the Mem Add key (numeric keypad))
	KEY_KP_MEMSUBTRACT     Key = sdl.K_KP_MEMSUBTRACT     // "Keypad MemSubtract" (the Mem Subtract key (numeric keypad))
	KEY_KP_MEMMULTIPLY     Key = sdl.K_KP_MEMMULTIPLY     // "Keypad MemMultiply" (the Mem Multiply key (numeric keypad))
	KEY_KP_MEMDIVIDE       Key = sdl.K_KP_MEMDIVIDE       // "Keypad MemDivide" (the Mem Divide key (numeric keypad))
	KEY_KP_PLUSMINUS       Key = sdl.K_KP_PLUSMINUS       // "Keypad +/-" (the +/- key (numeric keypad))
	KEY_KP_CLEAR           Key = sdl.K_KP_CLEAR           // "Keypad Clear" (the Clear key (numeric keypad))
	KEY_KP_CLEARENTRY      Key = sdl.K_KP_CLEARENTRY      // "Keypad ClearEntry" (the Clear Entry key (numeric keypad))
	KEY_KP_BINARY          Key = sdl.K_KP_BINARY          // "Keypad Binary" (the Binary key (numeric keypad))
	KEY_KP_OCTAL           Key = sdl.K_KP_OCTAL           // "Keypad Octal" (the Octal key (numeric keypad))
	KEY_KP_DECIMAL         Key = sdl.K_KP_DECIMAL         // "Keypad Decimal" (the Decimal key (numeric keypad))
	KEY_KP_HEXADECIMAL     Key = sdl.K_KP_HEXADECIMAL     // "Keypad Hexadecimal" (the Hexadecimal key (numeric keypad))

	KEY_LCTRL  Key = sdl.K_LCTRL  // "Left Ctrl"
	KEY_LSHIFT Key = sdl.K_LSHIFT // "Left Shift"
	KEY_LALT   Key = sdl.K_LALT   // "Left Alt" (alt, option)
	KEY_LGUI   Key = sdl.K_LGUI   // "Left GUI" (windows, command (apple), meta)
	KEY_RCTRL  Key = sdl.K_RCTRL  // "Right Ctrl"
	KEY_RSHIFT Key = sdl.K_RSHIFT // "Right Shift"
	KEY_RALT   Key = sdl.K_RALT   // "Right Alt" (alt, option)
	KEY_RGUI   Key = sdl.K_RGUI   // "Right GUI" (windows, command (apple), meta)

	KEY_MODE Key = sdl.K_MODE // "ModeSwitch" (I'm not sure if this is really not covered by any of the above, but since there's a special KMOD_MODE for it I'm adding it here)

	KEY_AUDIONEXT    Key = sdl.K_AUDIONEXT    // "AudioNext" (the Next Track media key)
	KEY_AUDIOPREV    Key = sdl.K_AUDIOPREV    // "AudioPrev" (the Previous Track media key)
	KEY_AUDIOSTOP    Key = sdl.K_AUDIOSTOP    // "AudioStop" (the Stop media key)
	KEY_AUDIOPLAY    Key = sdl.K_AUDIOPLAY    // "AudioPlay" (the Play media key)
	KEY_AUDIOMUTE    Key = sdl.K_AUDIOMUTE    // "AudioMute" (the Mute volume key)
	KEY_MEDIASELECT  Key = sdl.K_MEDIASELECT  // "MediaSelect" (the Media Select key)
	KEY_WWW          Key = sdl.K_WWW          // "WWW" (the WWW/World Wide Web key)
	KEY_MAIL         Key = sdl.K_MAIL         // "Mail" (the Mail/eMail key)
	KEY_CALCULATOR   Key = sdl.K_CALCULATOR   // "Calculator" (the Calculator key)
	KEY_COMPUTER     Key = sdl.K_COMPUTER     // "Computer" (the My Computer key)
	KEY_AC_SEARCH    Key = sdl.K_AC_SEARCH    // "AC Search" (the Search key (application control keypad))
	KEY_AC_HOME      Key = sdl.K_AC_HOME      // "AC Home" (the Home key (application control keypad))
	KEY_AC_BACK      Key = sdl.K_AC_BACK      // "AC Back" (the Back key (application control keypad))
	KEY_AC_FORWARD   Key = sdl.K_AC_FORWARD   // "AC Forward" (the Forward key (application control keypad))
	KEY_AC_STOP      Key = sdl.K_AC_STOP      // "AC Stop" (the Stop key (application control keypad))
	KEY_AC_REFRESH   Key = sdl.K_AC_REFRESH   // "AC Refresh" (the Refresh key (application control keypad))
	KEY_AC_BOOKMARKS Key = sdl.K_AC_BOOKMARKS // "AC Bookmarks" (the Bookmarks key (application control keypad))

	KEY_BRIGHTNESSDOWN Key = sdl.K_BRIGHTNESSDOWN // "BrightnessDown" (the Brightness Down key)
	KEY_BRIGHTNESSUP   Key = sdl.K_BRIGHTNESSUP   // "BrightnessUp" (the Brightness Up key)
	KEY_DISPLAYSWITCH  Key = sdl.K_DISPLAYSWITCH  // "DisplaySwitch" (display mirroring/dual display switch, video mode switch)
	KEY_KBDILLUMTOGGLE Key = sdl.K_KBDILLUMTOGGLE // "KBDIllumToggle" (the Keyboard Illumination Toggle key)
	KEY_KBDILLUMDOWN   Key = sdl.K_KBDILLUMDOWN   // "KBDIllumDown" (the Keyboard Illumination Down key)
	KEY_KBDILLUMUP     Key = sdl.K_KBDILLUMUP     // "KBDIllumUp" (the Keyboard Illumination Up key)
	KEY_EJECT          Key = sdl.K_EJECT          // "Eject" (the Eject key)
	KEY_SLEEP          Key = sdl.K_SLEEP          // "Sleep" (the Sleep key)
)
