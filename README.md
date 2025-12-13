# Scraper & Transform Component Guide

## Overview

This guide documents the Web Scraper component enhancements and Transform component for building dynamic web scraping workflows.

---

## Web Scraper Component (v0.0.8)

### New Features

#### 1. Form Submission (`submit_and_extract`)
Interact with web forms by filling fields and submitting.

**Parameters:**
- `form_data`: JSON object with CSS selectors and values
  - Example: `{"input[name='q']": "surabaya"}`
- `submit_selector`: CSS selector for submit button (optional)
  - If empty, presses Enter on last field
- `wait_time`: Milliseconds to wait after submission (default: 2000)

**Example:**
```json
{
  "action": "submit_and_extract",
  "method": "browser",
  "url": "https://example.com/search",
  "form_data": "{\"#searchBox\": \"weather\"}",
  "submit_selector": "button.submit",
  "wait_time": "5000",
  "extraction_rules": "{\"title\": \"h1\"}"
}
```

#### 2. Regex Extraction (`extract_with_regex`)
Extract data using regex patterns instead of CSS selectors.

**Parameters:**
- `extraction_rules`: JSON object with field names and regex patterns
- `wait_time`: Page load wait time

**Pattern Syntax:**
- Use `[0-9]` instead of `\d` to avoid escaping issues
- Use capture groups `()` to extract specific parts
- Without capture groups, returns full match

**Example:**
```json
{
  "action": "extract_with_regex",
  "url": "https://openweathermap.org/city/1625822",
  "extraction_rules": "{\"temperature\":\"([0-9]+)°C\",\"humidity\":\"Humidity:([0-9]+)%\"}"
}
```

**Output:**
```json
{
  "temperature": ["29", "27", "24"],
  "humidity": ["84"]
}
```

#### 3. Universal Wait Time
The `wait_time` parameter now works for ALL browser actions, not just form submission.

**Default:** 5000ms (5 seconds)

---

## Transform Component

Unified data transformation component with 3 modes.

### Mode 1: `array_to_first`
Converts all array values to their first element.

**Use Case:** Get current temperature instead of all forecast values

**Configuration:**
```json
{
  "transform_type": "array_to_first"
}
```

**Example:**
```
Input:  {"temperature": ["29", "27"], "humidity": ["84"]}
Output: {"temperature": "29", "humidity": "84"}
```

---

### Mode 2: `format_message`
Formats data into a message using template placeholders.

**Use Case:** Create readable weather reports, status messages

**Configuration:**
```json
{
  "transform_type": "format_message",
  "message_template": "Temperature: {{temperature}}°C\nHumidity: {{humidity}}%"
}
```

**Template Syntax:**
- Use `{{field_name}}` for placeholders
- Supports newlines with `\n`
- Arrays are automatically converted to first element

**Example:**
```
Input:  {"temperature": ["29"], "humidity": ["84"]}
Output: {
  "message": "Temperature: 29°C\nHumidity: 84%",
  "temperature": "29",
  "humidity": "84"
}
```

---

### Mode 3: `format_array`
Formats two parallel arrays into lines (key-value pairs).

**Use Case:** Currency rates, price lists, any paired data

**Configuration:**
```json
{
  "transform_type": "format_array",
  "key_field": "currency",
  "value_field": "rate",
  "separator": " = ",
  "message_template": "{key}{sep}{value}"
}
```

**Template Placeholders:**
- `{key}` - Value from key array
- `{value}` - Value from value array
- `{sep}` - Separator string

**Example:**
```
Input:  {"currency": ["USD", "EUR"], "rate": ["15000", "16000"]}
Output: {
  "message": "USD = 15000\nEUR = 16000",
  "count": 2
}
```

---

## Complete Weather Workflow Example

### Workflow: Dynamic City Weather

**Steps:**

1. **Search City** (Scraper)
   - URL: `https://openweathermap.org/find?q=$city_name`
   - Action: `extract_links`
   - Selector: `a[href*='/city/']`

2. **Get First Link** (Transform)
   - Transform Type: `array_to_first`

3. **Extract Weather** (Scraper)
   - URL: `$links`
   - Action: `extract_with_regex`
   - Rules: `{"temperature":"([0-9]+)°C","humidity":"Humidity:([0-9]+)%","wind":"([0-9]+.[0-9]+)m/s"}`

4. **Format Message** (Transform)
   - Transform Type: `format_message`
   - Template: `🌡️ {{temperature}}°C\n💧 {{humidity}}%\n💨 {{wind}} m/s`

5. **Send** (SendMessage)
   - Message: `$message`

---

## Important Notes

### CSS Selectors vs Regex

**Use CSS Selectors when:**
- Page has consistent HTML structure
- Elements have IDs or classes
- Need to extract from specific DOM elements

**Use Regex when:**
- Content is in plain text
- CSS selectors fail (dynamic classes)
- Need to extract patterns from text

### Variable Resolution

**In Transform:**
- Field names are LITERAL (no `$` prefix)
  - ✅ `key_field: "currency"`
  - ❌ `key_field: "$currency"`

**In SendMessage:**
- Use `$` to access variables from previous nodes
  - ✅ `message: "$message"`
  - ✅ `message: "Temp: $temperature"`

### Browser Automation

**Requirements:**
- Chrome/Chromium must be installed
- Use `method: "browser"` for JavaScript sites
- Increase `wait_time` for slow-loading pages

**Cross-Platform Paths:**
- Windows: `C:\Program Files\Google\Chrome\Application\chrome.exe`
- Linux: `/usr/bin/google-chrome`, `/usr/bin/chromium`
- macOS: `/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`

### Newlines in Messages

AI chat now supports newlines with `whitespace-pre-wrap` CSS.

**Usage:**
```
Temperature: 29°C
Humidity: 84%
```

Will display as 2 lines, not 1.

---

## Troubleshooting

### Scraper Returns Null

**Check:**
1. Is `method: "browser"` set for JavaScript sites?
2. Is `wait_time` long enough?
3. Are CSS selectors correct? (inspect element in browser)
4. For regex: are patterns escaped correctly?

### Transform Errors

**Common Issues:**
- Field names with `$` prefix (remove it)
- Wrong field names (check previous node output)
- Missing `message_template` for `format_message`

### Form Submission Fails

**Try:**
1. Use empty `submit_selector` to press Enter instead
2. Increase `wait_time`
3. Verify form field selectors in browser inspector

---

## Version History

- **v0.0.8** - Added `extract_with_regex`, universal `wait_time`
- **v0.0.6** - Renamed `wait_after_submit` to `wait_time`
- **v0.0.4** - Added Enter key submission support
- **v0.0.3** - Added form submission with `submit_and_extract`

---

## Quick Reference

### Scraper Actions
- `get_html` - Get raw HTML
- `extract_text` - Extract all text
- `extract_links` - Extract links
- `extract_images` - Extract images
- `extract_data` - Extract with CSS selectors
- `extract_one_data` - Extract single item
- `extract_with_regex` - Extract with regex patterns
- `submit_and_extract` - Submit form and extract

### Transform Types
- `array_to_first` - Get first element from arrays
- `format_message` - Format data into message
- `format_array` - Format parallel arrays into lines
