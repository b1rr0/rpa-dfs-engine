package traverser

// Common selector constants for web automation
var SelectorMap = map[string]string{
	// Form Elements
	"LOGIN_USERNAME": "#username",
	"LOGIN_PASSWORD": "#password", 
	"LOGIN_SUBMIT":   "#loginButton",
	"EMAIL_INPUT":    "input[type='email']",
	"SEARCH_INPUT":   "#searchInput",
	"SUBMIT_BUTTON":  "button[type='submit']",

	// Navigation Elements
	"HOME_LINK":   "a[href='/']",
	"BACK_BUTTON": ".back-button", 
	"NEXT_BUTTON": ".next-button",
	"MENU_TOGGLE": ".menu-toggle",

	// Common UI Elements
	"MODAL_CLOSE":     ".modal-close",
	"DROPDOWN_TOGGLE": ".dropdown-toggle",
	"CHECKBOX":        "input[type='checkbox']",
	"RADIO_BUTTON":    "input[type='radio']",
	"FILE_INPUT":      "input[type='file']",

	// Table Elements
	"TABLE_ROW":    "tr",
	"TABLE_CELL":   "td",
	"TABLE_HEADER": "th",

	// Status Elements
	"SUCCESS_MESSAGE": ".success-message",
	"ERROR_MESSAGE":   ".error-message", 
	"WARNING_MESSAGE": ".warning-message",
	"LOADING_SPINNER": ".loading-spinner",

	// Form Validation
	"REQUIRED_FIELD": "[required]",
	"INVALID_FIELD":  ".invalid",
	"VALID_FIELD":    ".valid",
}
