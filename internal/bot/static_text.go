package bot

import "path/filepath"

const API_KEY = "5360047799:AAE95DZ1rnPnxP5vLbkcOVREYGfPFqARbQs"

var BOT_USERNAME = ""
var ADMIN_WANTS_TO_GO_USER_MODE = false
var EMPTY_STRING = ""
var DELETE_STRING = "حذف"
var PICTURES_DIRECTORY = filepath.Join("../pictures/books/")

const (
	// Requests
	REQUEST_BOOK_ISBN               = "شابک کتاب را وارد کنید :"
	REQUEST_BOOK_TITLE              = "عنوان کتاب را وارد کنید :"
	REQUEST_BOOK_AUTHOR             = "نویسنده کتاب را وارد کنید :"
	REQUEST_BOOK_TRANSLATOR         = "مترجم کتاب را وارد کنید :"
	REQUEST_BOOK_PAPER_TYPE         = "نوع کاغذ کتاب را وارد کنید :"
	REQUEST_BOOK_WEIGHT             = "وزن کتاب را وارد کنید (گرم) :"
	REQUEST_BOOK_DESCRIPTION        = "توضیحاتی در مورد کتاب وارد کنید :"
	REQUEST_BOOK_NUMBER_OF_PAGES    = "تعداد صفحات کتاب را وارد کنید :"
	REQUEST_BOOK_GENRE              = "دسته بندی کتاب را وارد کنید :"
	REQUEST_BOOK_CENSORED_STATUS    = "وضعیت سانسور کتاب را انتخاب کنید :"
	REQUEST_BOOK_PUBLISHER          = "انتشارات کتاب را وارد کنید :"
	REQUEST_BOOK_PUBLISHDATE        = "تاریخ چاپ کتاب را وارد کنید :"
	REQUEST_BOOK_PRICE              = "قیمت کتاب را وارد کنید :"
	REQUEST_BOOK_AZERO_SCORE        = "امتیازی که ارزو به کتاب میدهد را وارد کنید :"
	REQUEST_BOOK_COVERTYPE          = "نوع جلد کتاب را انتخاب کنید :"
	REQUEST_BOOK_SIZE               = "سایز(قطع) کتاب را وارد کنید :"
	REQUEST_BOOK_AGE_CATEGORY       = "رده سنی کتاب را انتخاب کنید :"
	REQUEST_ADDRESS_COUNTRY         = "کدام کشور هستید؟"
	REQUEST_ADDRESS_PROVINCE        = "کدام استان هستید؟"
	REQUEST_ADDRESS_CITY            = "کدام شهر هستید؟"
	REQUEST_ADDRESS_STREET          = "کدام خیابان/کوچه هستید؟"
	REQUEST_ADDRESS_BUILDING_NUMBER = "کدام نام/شماره ساختمان هستید؟"
	REQUEST_ADDRESS_POSTAL_CODE     = "کد پستی شما چیست؟"
	REQUEST_ADDRESS_DESCRIPTION     = "میتوانید توضیحاتی را اینجا بنویسید!"
	REQUEST_ADDRESS_PHONE_NUMBER    = "شماره همراهتان را بنویسید :"
	REQUEST_BOOK_PICTURE            = "عکس های کتاب را در یک پیام ارسال کنید :"
	REQUEST_TRACKING_CODE           = "لطفا کد رهگیری سفارش را وارد کنید :"
	REQUEST_ORDER_REJECT_REASON     = "دلیل رد تایید سفارش را وارد کنید: "
	// Messages
	YOU_HAVE_NOT_ADDRESS_INLINE_KEYBOARD_MESSAGE      = "شما ادرسی را تنظیم نکرده اید, برای تنظیم ادرس بر روی دکمه ی زیر کلیک کنید !"
	CONTACT_TO_ADMIN_MESSAGE                          = "با کلیک بر روی دکمه ی زیر میتوانید با ادمین ارتباط داشته باشید !"
	START_TEXT                                        = "سلام خوش امدید"
	SHOW_ORDERS_HEADER_MESSAGE                        = "show order header"
	SHOW_ORDERS_FOOTER_MESSAGE                        = "اگر سوالی در رابطه با این سفارش دارید میتوانید سوالتان را همراه با کد رهگیری به ادمین ارسال کنید !"
	NO_ORDERS_IN_CONFIRMATION_QUEUE                   = "هیچ سفارشی در صف تایید نیست !"
	CONFIRMATION_OR_REJECTION_ORDRES_CANCELED_MESSAGE = "عملیات تایید یا رد تایید سفارشات لغو شد."
	ORDER_CONFIRMED_MESSAGE                           = "سفارش تایید شد و به کاربر اطلاع داده شد !"
	ORDER_REJECTED_MESSAGE                            = "سفارش رد شد و به کاربر اطلاع داده شد !"
	BUY_CART_MESSAGE_HEADER_MESSAGE                   = "لطفا مبلغ درج شده را به شماره کارت زیر واریز کنید و بعد <b>حتما</b> روی دکمه ی زیر کلیک کنید.\nپس از ان سفارش شما در صف تایید قرار میگیرد و میتوانید از طریق منو اصلی برنامه بر روی دکمه ی سفارش ها کلیک کنید و وضعیت سفارش خود را ببینید.\n"
	BUY_CART_MESSAGE_FOOTER_MESSAGE                   = "buy cart footer"
	CART_MESSAGE_HEADER_MESSAGE                       = "شما میتوانید سبد خریدتان را در اینجا مشاهده کنید !"
	CART_MESSAGE_FOOTER_MESSAGE                       = "با کلیک بر روی دکمه ی زیر میتوانید سبد خریدتان را بخرید !"
	FAQ_MESSAGE                                       = "شما میتوانید سوالات متداولی که پرسیده میشود را در متن زیر بخوانید !\n"
	ORDER_ADDED_FORMAT                                = "سفارش شما ثبت شد !\nشما میتوانید از طریق منوی اصلی وضعیت سفارش خود را مشاهده کنید.\nهمچنین در صورت نیاز میتوانید سوالات خود را از ادمین بپرسید.\n@%s"
	ADMIN_START_TEXT                                  = "سلام ادمین خوش اومدی"
	NO_RESULT_FOUND_DESCRIPTION_FORMAT_ERROR          = "برای عبارت %s نتیجه ای یافت نشد.\nعنوان کتاب را بررسی کنید, همچنین امکان دارد کتاب موجود نباشد."
	SEARCH_TEXT                                       = "با کلیک بر روی دکمه ی زیر میتوانید کتاب مورد نظر خود را جستجو کنید.\nحتما نام کتاب را به فارسی جستجو کنید.\n"
	YOU_ALREADY_HAVE_ADDRESS                          = "شما قبلا ادرسی را وارد کرده اید !"
	SHOW_USER_ADDRESS_FORMATTED                       = "کشور: %s\nاستان: %s\nشهر: %s\nخیابان: %s\nشماره ی ساختمان: %s\nکد پستی: %s\nشماره تلفن: %s\nتوضیحات: %s\n"
	BOOK_INFORMATION_FORMAT                           = "عنوان: %s\nنویسنده: %s\nمترجم: %s\nتعداد صفحات: %d\nدسته بندی: %s\nوضعیت سانسور: %s\nنوع جلد: %s\nسایز: %s\nرده سنی: %s\nامتیاز گودریدز: %s\nامتیاز ارزو: %s\nانتشارات: %s\nتاریخ انتشار: %s\nشابک: %s\nقیمت: %s\n\n@%s"
	BOOK_DELETED_MESSAGE                              = "کتاب %s حدف شد."
	BOOK_DELETED_FROM_CART_FORMAT                     = "کتاب %s از سبد خرید حذف شد."
	BOOK_NOT_ADDED_TO_CART_MESSAGE                    = "کتاب به سبد حرید اضافه نشد.\nدوباره امتحان کنید."
	BOOK_ADDED_TO_CART_MESSAGE                        = "کتاب به سبد خرید اضافه شد."
	THERE_IS_NO_IN_CONFIRMATION_ORDER_LEFT            = "سفارش دیگری در صف تایید باقی نمانده است !"
	SWITCH_TO_PV_FORMAT                               = "https://t.me/%s"
	ADDRESS_ADDED_MESSAGE                             = "ادرس شما تنظیم شد."
	ADD_BOOK_FROM_CART_MESSAGE                        = "با کلیک بر روی دکمه ی زیر میتوانید کتاب را به سبد خرید اضافه کنید !"
	REMOVE_BOOK_FROM_CART_MESSAGE                     = "با کلیک بر روی دکمه ی زیر میتوانید کتاب را از سبد خرید حذف کنید !"
	// Errors
	UNKNOWN_ERROR                        = "مشکلی پیش امد, دوباره امتحان کنید."
	ENTERED_VALUE_IS_INVALID_ERROR       = "مقدار وارد شده نامعتبر است."
	ENTERED_NON_NUMBER_VALUE_ERROR       = "لطفا عددی را به عنوان مقدار وارد کنید."
	BOOK_NOT_SAVED_IN_DATABASE_ERROR     = "عملیات ذخیره کردن کتاب در دیتابیس ناموفق بود."
	AT_LEAST_ENTER_ONE_CHARACTER_ERROR   = "لطفا حداقل یک کاراکتر برای جستجو وارد کنید."
	ENTERED_PHRASE_IS_TOO_SHORT_ERROR    = "عبارت وارد شده بیش از حد کوتاه است."
	NO_RESULT_FOUND_ERROR                = "نتیجه ای یافت نشد."
	LENGTH_OF_TRACKING_CODE_IS_INCORRECT = "طول کد رهگیری وارد شده اشتباه است !"

	// Keyboards item title
	SEARCH_BOOK_KEYBOARD_ITEM_TITLE            = "جستجوی کتاب"
	CART_KEYBOARD_ITEM_TITLE                   = "سبد خرید"
	BUY_CART_KEYBOARD_ITEM_TITLE               = "خرید سبد"
	CONTACT_ADMIN_KEYBOARD_ITEM_TITLE          = "ارتباط با ادمین"
	ADDRESS_KEYBOARD_ITEM_TITLE                = "تنظیم ادرس"
	SHOW_ORDERS_KEYBOARD_ITEM_TITLE            = "سفارشات"
	FAQ_KEYBOARD_ITEM_TITLE                    = "سوالات متداول"
	ADD_POST_TRACKING_CODE_KEYBOARD_ITEM_TITLE = "افزودن کد رهگیری"

	MAIN_MENU_ITEM_TITLE       = "منو اصلی"
	CANCEL_KEYBOARD_ITEM_TITLE = "انصراف"

	ADMIN_ADD_BOOK_KEYBOARD_ITEM_TITLE         = "اضافه کردن کتاب"
	ADMIN_DELETE_BOOK_KEYBOARD_ITEM_TITLE      = "حدف کتاب"
	ADMIN_CONFIRM_ORDERS_KEYBOARD_ITEM_TITLE   = "تأیید سفارشات"
	ADMIN_STATISTICS_KEYBOARD_ITEM_TITLE       = "امار ربات"
	ADMIN_BACK_TO_USER_PANEL_ITEM_TITLE        = "پنل کاربر"
	ADMIN_BACK_TO_ADMIN_PANEL_ITEM_TITLE       = "پنل ادمین"
	ADMIN_CHECK_ORDER_BY_TRACKING_CODE         = "بررسی سفارش"
	IS_CENSORED_STATUS_KEYBOARD_ITEM_TITLE     = "سانسور شده است"
	IS_NOT_CENSORED_STATUS_KEYBOARD_ITEM_TITLE = "بدون سانسور است"

	ENTER_SEARCH_PHRASE_TEXT                 = "وارد کردن عبارت حستجو"
	ENTER_SEARCH_PHRASE_FOR_DELETE_BOOK_TEXT = "جستجو برای حذف"
	SEARCH_FOR_DELETE_BOOK_MESSAGE_TEXT      = "با استفاده از گزینه ی زیر میتوانید کتابی را برای حذف جستجو کنید."

	// Inline keyboards item title
	ADD_BOOK_TO_CART_INLINE_KEYBOARD_ITEM_TITLE       = "اضافه کردن به سبد خرید"
	DELETE_BOOK_FROM_CART_INLINE_KEYBOARD_ITEM_TITLE  = "حدف از سبد خرید"
	CLICK_FOR_EDIT_ADDRESS_INLINE_KEYBOARD_ITEM_TITLE = "برای تغییر ادرس کلیک کنید !"
	I_PAID_CART_INLINE_KEYBOARD_ITEM_TITLE            = "پرداخت کردم"
	CONFIRM_ORDER_KEYBOARD_ITEM                       = "تایید سفارش"
	REJECT_ORDER_KEYBOARD_ITEM                        = "رد سفارش"

	// Order
	SEND_ORDER_CONFIRMED_TO_USER_HEADER = "سفارش شما تایید شد و در حال بسته بندی و نهایتا ارسال است"
	SEND_ORDER_CONFIRMED_TO_USER_FOOTER = "اتمام پیام !"
	SEND_ORDER_REJECTED_TO_USER_HEADER  = "سفارش شما رد شد !"
	SEND_ORDER_REJECTED_TO_USER_FOOTER  = "اتمام پیام !"
	// Miscs
	BOT_START_QUERY = "https://t.me/%s/?start=%d"
	CART_IS_EMPTY   = "سبد خرید خالی هست !"
)
