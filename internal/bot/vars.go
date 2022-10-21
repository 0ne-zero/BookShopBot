package bot

import "path/filepath"

const API_KEY = "5360047799:AAE95DZ1rnPnxP5vLbkcOVREYGfPFqARbQs"

var BOT_USERNAME = ""

const UNKNOWN_ERROR = "مشکلی پیش امد, دوباره امتحان کنید."

var PICTURES_DIRECTORY = filepath.Join("../pictures/books/")

var BOT_START_QUERY = "https://t.me/%s/?start=%d"

var ENTERED_VALUE_IS_INVALID_ERROR = "مقدار وارد شده نامعتبر است."
var BOOK_NOT_SAVED_IN_DATABASE = "عملیات ذخیره کردن کتاب در دیتابیس ناموفق بود."
var ENTERED_NON_NUMBER_VALUE_ERROR = "لطفا عددی را به عنوان مقدار وارد کنید."
var EMPTY_STRING = ""
var DELETE_STRING = "حذف"
var ENTER_SEARCH_PHRASE_TEXT = "وارد کردن عبارت حستجو"
var ENTER_SEARCH_PHRASE_FOR_DELETE_BOOK_TEXT = "وارد کردن عبارت جستجو برای حذف"
var SEARCH_FOR_DELETE_BOOK_MESSAGE_TEXT = "با استفاده از گزینه ی زیر میتوانید کتابی را برای حذف انتخاب کنید"
var START_TEXT = "start"
var ADMIN_START_TEXT = "admin start"
var AT_LEAST_ENTER_ONE_CHARACTER_ERROR = "لطفا حداقل یک کاراکتر برای جستجو وارد کنید."
var ENTERED_PHRASE_IS_TOO_SHORT_ERROR = "عبارت وارد شده بیش از حد کوتاه است."
var NO_RESULT_FOUND_ERROR = "نتیجه ای یافت نشد."
var NO_RESULT_FOUND_DESCRIPTION_ERROR = "برای عبارت %s نتیجه ای یافت نشد.\nعنوان کتاب را بررسی کنید, همچنین امکان دارد کتاب موجود نباشد."
var SEARCH_TEXT = "you can search by pressing below button"
var REQUEST_BOOK_PICTURE = "عکس های کتاب را در یک پیام ارسال کنید."
var REQUEST_BOOK_ISBN = "isbn"
var REQUEST_BOOK_TITLE = "عنوان کتاب را وارد کنید :"
var REQUEST_BOOK_AUTHOR = "author"
var REQUEST_BOOK_TRANSLATOR = "translator"
var REQUEST_BOOK_PAPER_TYPE = "paper type"
var REQUEST_BOOK_DESCRIPTION = "description"
var REQUEST_BOOK_NUMBER_OF_PAGES = "number of pages"
var REQUEST_BOOK_GENRE = "genre"
var REQUEST_BOOK_CENSORED_STATUS = "censored status"
var REQUEST_BOOK_PUBLISHER = "publisher"
var REQUEST_BOOK_PUBLISHDATE = "publish date"
var REQUEST_BOOK_PRICE = "price"
var REQUEST_BOOK_AZERO_SCORE = "arezo score"
var REQUEST_BOOK_COVERTYPE = "select cover type"
var REQUEST_BOOK_SIZE = "select book size"
var REQUEST_BOOK_AGE_CATEGORY = "select age category"

var BOOK_INFORMATION_FORMAT = "عنوان: %s\nنویسنده: %s\nمترجم: %s\nتعداد صفحات: %d\nدسته بندی: %s\nوضعیت سانسور: %s\nنوع جلد: %s\nسایز: %s\nرده سنی: %s\nامتیاز گودریدز: %s\nامتیاز ارزو: %s\nانتشارات: %s\nتاریخ انتشار: %s\nشابک: %s\nقیمت: %s\n\n@%s"
var BOOK_DELETED_SUCCESSFULY = "کتاب %s حدف شد."

var ADD_BOOK_TO_CART = "اضافه کردن به سبد خرید"
var BOOK_NOT_ADDED_TO_CART_SUCCESSFULY = "کتاب به سبد حرید اضافه نشد."
var BOOK_ADDED_TO_CART_SUCCESSFULY = "کتاب به سبد خرید اضافه شد."
