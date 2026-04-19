package categorizer

// categoryMapping pairs a top-level category with a subcategory.
type categoryMapping struct {
	Category    string
	SubCategory string
}

// mccCodes maps Merchant Category Codes (ISO 18245) to spend categories.
// Used by GetCategories to enumerate all possible categories.
var mccCodes = map[string]categoryMapping{
	// Grocery
	"5411": {"Grocery", "Supermarket"},
	"5422": {"Grocery", "Fresh Produce"},
	"5441": {"Grocery", "Candy & Confectionery"},
	"5451": {"Grocery", "Dairy"},
	"5462": {"Grocery", "Bakery"},
	"5499": {"Grocery", "General"},

	// Food & Dining
	"5812": {"Food & Dining", "Restaurants"},
	"5813": {"Food & Dining", "Bars & Lounges"},
	"5814": {"Food & Dining", "Fast Food"},

	// Travel
	"3000": {"Travel", "Airlines"},
	"3001": {"Travel", "Airlines"},
	"4511": {"Travel", "Airlines"},
	"4722": {"Travel", "Travel Agency"},
	"7011": {"Travel", "Hotels"},
	"7012": {"Travel", "Hotels"},
	"3501": {"Travel", "Hotels"},
	"7512": {"Travel", "Car Rental"},
	"4121": {"Travel", "Ride Sharing"},
	"4111": {"Travel", "Public Transport"},
	"4112": {"Travel", "Railways"},
	"4131": {"Travel", "Bus"},
	"4411": {"Travel", "Cruise"},
	"7523": {"Travel", "Parking"},
	"5541": {"Travel", "Fuel"},
	"5542": {"Travel", "Fuel"},

	// Shopping
	"5651": {"Shopping", "Clothing"},
	"5611": {"Shopping", "Clothing"},
	"5621": {"Shopping", "Clothing"},
	"5631": {"Shopping", "Clothing"},
	"5641": {"Shopping", "Clothing"},
	"5661": {"Shopping", "Footwear"},
	"5691": {"Shopping", "Clothing"},
	"5699": {"Shopping", "Clothing"},
	"5944": {"Shopping", "Jewellery"},
	"5945": {"Shopping", "Toys & Games"},
	"5946": {"Shopping", "Photography"},
	"5947": {"Shopping", "Gifts"},
	"5732": {"Shopping", "Electronics"},
	"5734": {"Shopping", "Software"},
	"5735": {"Shopping", "Music"},
	"5045": {"Shopping", "Electronics"},
	"5065": {"Shopping", "Electronics"},
	"5200": {"Shopping", "Home & Garden"},
	"5211": {"Shopping", "Home & Garden"},
	"5231": {"Shopping", "Home & Garden"},
	"5251": {"Shopping", "Hardware"},
	"5261": {"Shopping", "Garden"},
	"5311": {"Shopping", "Department Store"},
	"5310": {"Shopping", "Discount Store"},
	"5331": {"Shopping", "Variety Store"},
	"5399": {"Shopping", "General Merchandise"},
	"5942": {"Shopping", "Books"},
	"5943": {"Shopping", "Stationery"},
	"5262": {"Shopping", "Marketplace"},
	"5964": {"Shopping", "Online Shopping"},
	"5965": {"Shopping", "Online Shopping"},
	"5966": {"Shopping", "Online Shopping"},

	// Insurance
	"6300": {"Insurance", "General"},
	"6381": {"Insurance", "General"},
	"5960": {"Insurance", "Health"},

	// Entertainment
	"7832": {"Entertainment", "Movies"},
	"7829": {"Entertainment", "Movies"},
	"7841": {"Entertainment", "Streaming"},
	"7911": {"Entertainment", "Events"},
	"7922": {"Entertainment", "Events"},
	"7929": {"Entertainment", "Events"},
	"7932": {"Entertainment", "Sports"},
	"7933": {"Entertainment", "Sports"},
	"7941": {"Entertainment", "Sports"},
	"7991": {"Entertainment", "Attractions"},
	"7992": {"Entertainment", "Attractions"},
	"7993": {"Entertainment", "Gaming"},
	"7994": {"Entertainment", "Gaming"},
	"7995": {"Entertainment", "Gaming"},
	"7996": {"Entertainment", "Attractions"},
	"7997": {"Entertainment", "Recreation"},
	"7998": {"Entertainment", "Aquariums"},
	"7999": {"Entertainment", "Recreation"},

	// Utilities
	"4814": {"Utilities", "Phone"},
	"4812": {"Utilities", "Phone"},
	"4813": {"Utilities", "Phone"},
	"4816": {"Utilities", "Internet"},
	"4899": {"Utilities", "Cable & TV"},
	"4900": {"Utilities", "Electricity"},
	"4901": {"Utilities", "Gas"},

	// Healthcare
	"8011": {"Healthcare", "Medical"},
	"8021": {"Healthcare", "Dental"},
	"8031": {"Healthcare", "Medical"},
	"8041": {"Healthcare", "Eye Care"},
	"8042": {"Healthcare", "Eye Care"},
	"8043": {"Healthcare", "Eye Care"},
	"8049": {"Healthcare", "Medical"},
	"8050": {"Healthcare", "Medical"},
	"8062": {"Healthcare", "Hospital"},
	"8071": {"Healthcare", "Medical Lab"},
	"8099": {"Healthcare", "Medical"},
	"5912": {"Healthcare", "Pharmacy"},
	"5122": {"Healthcare", "Pharmacy"},

	// Education
	"8211": {"Education", "Schools"},
	"8220": {"Education", "University"},
	"8241": {"Education", "Courses"},
	"8244": {"Education", "Courses"},
	"8249": {"Education", "Courses"},
	"8299": {"Education", "Courses"},
	"5111": {"Education", "Stationery"},
	"5192": {"Education", "Books"},

	// Lifestyle
	"7032": {"Lifestyle", "Fitness"},
	"7298": {"Lifestyle", "Spa & Beauty"},
	"7230": {"Lifestyle", "Beauty Salon"},
	"7297": {"Lifestyle", "Spa & Beauty"},
	"5977": {"Lifestyle", "Cosmetics"},

	// Financial
	"6010": {"Financial", "Cash Withdrawal"},
	"6011": {"Financial", "Cash Withdrawal"},
	"6012": {"Financial", "Bank Fees"},
	"6051": {"Financial", "Transfer"},
	"6211": {"Financial", "Securities"},
	"6513": {"Financial", "Rent"},
	"9211": {"Financial", "Government Fees"},
	"9222": {"Financial", "Government Fees"},
	"9311": {"Financial", "Tax"},
	"9399": {"Financial", "Government"},
	"9402": {"Financial", "Postal"},
}

// keywordRule maps one or more substrings to a category. Rules are evaluated
// in order; the first matching keyword wins. More specific rules (e.g. merchant
// names) should appear before generic ones (e.g. "restaurant").
type keywordRule struct {
	keywords    []string
	category    string
	subCategory string
}

// keywordRules is the ordered list of keyword-to-category mappings.
var keywordRules = []keywordRule{
	// Grocery
	{[]string{"bigbasket", "big basket", "blinkit", "zepto", "dunzo", "grofers", "jiomart", "dmart", "d-mart", "reliance fresh", "more supermarket", "nature basket", "spencers", "star bazaar", "walmart", "costco", "kroger", "whole foods", "trader joe", "aldi", "lidl", "tesco", "sainsbury", "woolworths", "super market", "supermarket", "ragam", "midas daily", "provisions"}, "Grocery", "Supermarket"},
	{[]string{"fresh", "organic", "vegetable", "fruit"}, "Grocery", "Fresh Produce"},
	{[]string{"sweets", "sweet shop", "haldiram", "bikanervala", "anand sweets", "kanti sweets", "gowri gajanand", "mithai", "namkeen"}, "Grocery", "Sweets & Bakery"},

	// Food & Dining
	{[]string{"swiggy", "rsp*swiggy", "raz*swiggy", "bundl technologies", "zomato", "uber eats", "doordash", "grubhub", "deliveroo", "foodpanda"}, "Food & Dining", "Delivery"},
	{[]string{"starbucks", "costa coffee", "cafe coffee day", "ccd", "blue tokai", "third wave", "tim horton", "dunkin", "chaayos"}, "Food & Dining", "Cafes"},
	{[]string{"mcdonald", "burger king", "kfc", "domino", "pizza hut", "subway", "wendy", "taco bell", "chipotle", "pizza bakery"}, "Food & Dining", "Fast Food"},
	{[]string{"eazydiner", "dineout", "restaurant", "bistro", "diner", "kitchen", "grill", "eatery", "dhaba", "mess", "jalpan", "muhavra"}, "Food & Dining", "Restaurants"},

	// Travel
	{[]string{"uber", "ola", "lyft", "rapido", "grab"}, "Travel", "Ride Sharing"},
	{[]string{"indigo", "air india", "spicejet", "vistara", "goair", "emirates", "british airways", "lufthansa", "singapore airlines", "qatar airways", "delta", "united airlines", "american airlines", "southwest", "ryanair", "easyjet", "jetblue"}, "Travel", "Airlines"},
	{[]string{"marriott", "hilton", "hyatt", "taj", "oberoi", "itc hotel", "radisson", "holiday inn", "airbnb", "oyo", "booking.com", "trivago", "makemytrip", "goibibo", "cleartrip", "expedia", "agoda"}, "Travel", "Hotels"},
	{[]string{"irctc", "railway", "metro", "bus", "transport"}, "Travel", "Public Transport"},
	{[]string{"shell", "hp petrol", "bharat petroleum", "indian oil", "bpcl", "hpcl", "iocl", "petrol", "diesel", "fuel"}, "Travel", "Fuel"},
	{[]string{"parking"}, "Travel", "Parking"},

	// Shopping - Online
	{[]string{"amazon", "flipkart", "myntra", "smartbuy", "gyftr", "ajio", "nykaa", "meesho", "snapdeal", "shopee", "lazada", "ebay", "etsy", "aliexpress", "shein", "asos", "zalando"}, "Shopping", "Online Shopping"},
	// Shopping - Electronics
	{[]string{"samsung", "croma", "reliance digital", "best buy", "currys", "mediamarkt"}, "Shopping", "Electronics"},
	// Shopping - Clothing & Fashion
	{[]string{"h&m", "zara", "uniqlo", "gap", "levi", "nike", "adidas", "puma", "reebok", "decathlon", "pantaloons", "westside", "life style", "lifestyle international", "shoppers stop", "aditya birla fashion", "vero moda", "allen solly", "van heusen", "peter england", "louis philippe", "jack & jones", "forever 21", "mango", "marks & spencer"}, "Shopping", "Clothing"},
	{[]string{"ikea", "home centre", "home depot", "pottery barn"}, "Shopping", "Home & Garden"},
	{[]string{"island star mall", "mall develo", "heera novelty"}, "Shopping", "General Merchandise"},
	{[]string{"stationery", "raj stationery"}, "Shopping", "Stationery"},

	// Insurance
	{[]string{"lic", "max life", "hdfc life", "icici prudential", "sbi life", "bajaj allianz", "star health", "new india assurance", "insurance", "premium", "policy bazaar", "policybazaar"}, "Insurance", "General"},

	// Entertainment
	{[]string{"netflix", "disney", "hotstar", "prime video", "hbo", "hulu", "spotify", "apple music", "youtube premium", "sony liv", "zee5", "voot", "jiocinema"}, "Entertainment", "Streaming"},
	{[]string{"pvr", "inox", "cinepolis", "imax", "bookmyshow", "cinema", "movie", "amc"}, "Entertainment", "Movies"},
	{[]string{"playstation", "xbox", "steam", "epic games", "nintendo", "gaming"}, "Entertainment", "Gaming"},
	{[]string{"concert", "event", "ticketmaster", "festival"}, "Entertainment", "Events"},
	{[]string{"playo", "sportsvilla"}, "Entertainment", "Sports & Recreation"},

	// Utilities
	{[]string{"jio", "airtel", "vodafone", "bsnl", "vi ", "t-mobile", "verizon", "at&t", "ee ", "o2 ", "three"}, "Utilities", "Phone"},
	{[]string{"broadband", "fiber", "wifi", "act fibernet", "hathway", "comcast", "spectrum"}, "Utilities", "Internet"},
	{[]string{"electricity", "power", "adani", "tata power", "bescom", "bses"}, "Utilities", "Electricity"},
	{[]string{"gas bill", "piped gas", "mahanagar gas", "indraprastha gas"}, "Utilities", "Gas"},
	{[]string{"water bill", "water supply"}, "Utilities", "Water"},

	// Healthcare
	{[]string{"hospital", "clinic", "diagnostic", "pathology"}, "Healthcare", "Medical"},
	{[]string{"pharmacy", "chemist", "medicine", "drug store", "apollo pharmacy", "medplus", "netmeds", "pharmeasy", "1mg", "walgreens", "cvs", "ranuja medical"}, "Healthcare", "Pharmacy"},
	{[]string{"dental", "dentist", "orthodont"}, "Healthcare", "Dental"},
	{[]string{"optician", "eye care", "lenskart", "vision", "optical"}, "Healthcare", "Eye Care"},

	// Education
	{[]string{"school", "academy", "college", "university", "institute", "coursera", "udemy", "edx", "skillshare", "masterclass", "linkedin learning", "pluralsight", "byju", "unacademy", "vedantu"}, "Education", "Courses"},
	{[]string{"book", "kindle"}, "Education", "Books"},

	// Lifestyle
	{[]string{"gym", "fitness", "cult.fit", "cultfit", "gold gym", "anytime fitness", "yoga", "crossfit"}, "Lifestyle", "Fitness"},
	{[]string{"salon", "spa", "beauty", "parlour", "parlor", "haircut", "grooming", "urban company", "urbanclap", "unisex salon"}, "Lifestyle", "Spa & Beauty"},

	// Software & Subscriptions
	{[]string{"defmacro", "cleartax", "notion", "slack", "github", "atlassian", "adobe", "microsoft", "google workspace", "dropbox", "canva"}, "Shopping", "Software & Subscriptions"},

	// Financial
	{[]string{"imps", "neft", "rtgs", "upi"}, "Financial", "Transfer"},
	{[]string{"emi", "loan", "installment"}, "Financial", "EMI"},
	{[]string{"interest charge", "finance charge", "late fee", "annual fee", "service charge", "redemption proc fee"}, "Financial", "Bank Fees"},
	{[]string{"igst", "cgst", "sgst", "gst-", "dcc transaction"}, "Financial", "Taxes & Fees"},
	{[]string{"rent", "pg ", "paying guest", "housing"}, "Financial", "Rent"},
	{[]string{"atm", "cash withdrawal"}, "Financial", "Cash Withdrawal"},
	{[]string{"mutual fund", "sip ", "investment", "zerodha", "groww", "kuvera", "paytm money"}, "Financial", "Investment"},
}
