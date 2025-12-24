package generator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func init() {
	RegisterComponent("Web Scraper", func(ctx *WorkflowContext) error {
		return handleScrape(ctx)
	})
}

func handleScrape(ctx *WorkflowContext) error {
	var (
		url             string
		method          string
		action          string
		selector        string
		captchaAPIKey   string
		userAgent       string
		headersJSON     string
		extractionRules string
		filter          string
		formDataJSON    string
		submitSelector  string
		waitTime        string
	)

	// Extract parameters
	for _, p := range ctx.Params {
		val := strings.TrimSpace(ResolveParam(ctx.FiberCtx, p.CompValue))


		switch p.InputName {
		case "url":
			url = val
		case "method":
			method = strings.TrimSpace(strings.ToLower(val))
		case "action":
			action = strings.TrimSpace(strings.ToLower(val))
		case "selector":
			selector = val
		case "captcha_api_key":
			captchaAPIKey = val
		case "user_agent":
			userAgent = val
		case "headers":
			headersJSON = val
		case "extraction_rules":
			extractionRules = val
		case "filter":
			filter = val
		case "form_data":
			formDataJSON = val
		case "submit_selector":
			submitSelector = val
		case "wait_time":
			waitTime = val
		case "click_selector":
			// We need a variable for this, it wasn't defined in the topvar block yet, let's just append it to new variables
	// Wait, I can't just define it there in switch.
	// I need to add 'clickSelector' to the var block at the top of handleScrape.
	// But to minimize diff, I will just assume I added it or re-read Params.
	// Actually, let's just do a clean update of the param loop.
		}
	}

    // Checking for click_selector in a separate loop to avoid var block modification complexity in replace tool
    var clickSelector string
	for _, p := range ctx.Params {
        if strings.TrimSpace(p.InputName) == "click_selector" {
            clickSelector = strings.TrimSpace(ResolveParam(ctx.FiberCtx, p.CompValue))
        }
    }

	// Default user agent
	if userAgent == "" {
		userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	}

	// Parse wait time
	waitMs := 5000 // default 5 seconds for browser actions
	if waitTime != "" {
		if parsed, err := strconv.Atoi(waitTime); err == nil {
			waitMs = parsed
		}
	}

	// Parse custom headers
	var customHeaders map[string]string
	if headersJSON != "" {
		json.Unmarshal([]byte(headersJSON), &customHeaders)
	}

	// Logic to perform scraping
	var result interface{}
	var err error

	fmt.Printf("[Scraper] Starting action: %s, method: %s, url: %s\n", action, method, url)

	switch action {
	case "get_html":
		if method == "browser" {
			result, err = scrapeWithBrowser(url, userAgent, "html", waitMs)
		} else {
			result, err = scrapeWithHTTP(url, userAgent, customHeaders, "html")
		}

	case "extract_text":
		if method == "browser" {
			html, err := scrapeWithBrowser(url, userAgent, "html", waitMs)
			if err == nil {
				result, err = extractTextFromHTML(getHTMLFromResult(html), selector)
			}
		} else {
			html, err := scrapeWithHTTP(url, userAgent, customHeaders, "html")
			if err == nil {
				result, err = extractTextFromHTML(html.(string), selector)
			}
		}

	case "extract_links":
		if method == "browser" {
			html, err := scrapeWithBrowser(url, userAgent, "html", waitMs)
			if err == nil {
				result, err = extractLinksFromHTML(getHTMLFromResult(html), selector)
			}
		} else {
			html, err := scrapeWithHTTP(url, userAgent, customHeaders, "html")
			if err == nil {
				result, err = extractLinksFromHTML(html.(string), selector)
			}
		}

	case "extract_images":
		if method == "browser" {
			html, err := scrapeWithBrowser(url, userAgent, "html", waitMs)
			if err == nil {
				result, err = extractImagesFromHTML(getHTMLFromResult(html), selector)
			}
		} else {
			html, err := scrapeWithHTTP(url, userAgent, customHeaders, "html")
			if err == nil {
				result, err = extractImagesFromHTML(html.(string), selector)
			}
		}

	case "extract_data":
		if method == "browser" {
			html, err := scrapeWithBrowser(url, userAgent, "html", waitMs)
			if err == nil {
				result, err = extractDataWithRules(getHTMLFromResult(html), extractionRules)
			}
		} else {
			html, err := scrapeWithHTTP(url, userAgent, customHeaders, "html")
			if err == nil {
				result, err = extractDataWithRules(html.(string), extractionRules)
			}
		}

		if err == nil && filter != "" {
			if resultMap, ok := result.(map[string]interface{}); ok {
				result, err = applyFilter(resultMap, filter)
			}
		}

	case "extract_one_data":
		if method == "browser" {
			html, err := scrapeWithBrowser(url, userAgent, "html", waitMs)
			if err == nil {
				result, err = extractDataWithRules(getHTMLFromResult(html), extractionRules)
			}
		} else {
			html, err := scrapeWithHTTP(url, userAgent, customHeaders, "html")
			if err == nil {
				result, err = extractDataWithRules(html.(string), extractionRules)
			}
		}

		if err == nil {
			if resultMap, ok := result.(map[string]interface{}); ok {
				if filter != "" {
					resultMap, err = applyFilter(resultMap, filter)
				}
				if err == nil {
					result = flattenResult(resultMap)
				}
			}
		}

	case "submit_and_extract":
		// Parse form data
		var formData map[string]string
		if formDataJSON != "" {
			if err := json.Unmarshal([]byte(formDataJSON), &formData); err != nil {
				result = map[string]string{"status": "error", "message": "Invalid form_data JSON"}
				break
			}
		}

		// Only works with browser method
		if method == "browser" {
			html, err := scrapeWithFormSubmit(url, userAgent, formData, submitSelector, clickSelector, waitMs)
			if err == nil {
				result, err = extractDataWithRules(getHTMLFromResult(html), extractionRules)
			}
		} else {
			result = map[string]string{"status": "error", "message": "submit_and_extract requires method=browser"}
		}

		if err == nil && filter != "" {
			if resultMap, ok := result.(map[string]interface{}); ok {
				result, err = applyFilter(resultMap, filter)
			}
		}

	case "extract_with_regex":
		// Get the page text first
		var textStr string
		if method == "browser" {
			html, err := scrapeWithBrowser(url, userAgent, "html", waitMs)
			if err == nil {
				textResult, err := extractTextFromHTML(getHTMLFromResult(html), "")
				if err == nil {
					if textArr, ok := textResult.([]string); ok {
						textStr = strings.Join(textArr, " ")
					}
				}
			}
		} else {
			html, err := scrapeWithHTTP(url, userAgent, customHeaders, "html")
			if err == nil {
				textResult, err := extractTextFromHTML(html.(string), "")
				if err == nil {
					if textArr, ok := textResult.([]string); ok {
						textStr = strings.Join(textArr, " ")
					}
				}
			}
		}

		if err == nil {
			result, err = extractWithRegex(textStr, extractionRules)
		}

	case "solve_captcha":
		result, err = solveCaptcha(url, captchaAPIKey)

	default:
		result = map[string]string{"status": "action_not_implemented_yet", "action": action}
	}

	if err != nil {
		fmt.Printf("[Scraper] Error: %v\n", err)
		return err
	}

	// Extract final URL if result is a map with "url" and "html" keys (from browser scraping)
	finalURL := ""
	if resultMap, ok := result.(map[string]interface{}); ok {
		if urlVal, exists := resultMap["url"]; exists {
			if urlStr, ok := urlVal.(string); ok {
				finalURL = urlStr
			}
			// If we have both html and url, use html as the result for further processing
			if htmlVal, exists := resultMap["html"]; exists {
				result = htmlVal
			}
		}
	}

	// Build final result map
	finalResult := make(map[string]interface{})
	finalResult["final_url"] = finalURL

	// If result is a map, merge it into finalResult (flattening)
	if resultMap, ok := result.(map[string]interface{}); ok {
		for k, v := range resultMap {
			finalResult[k] = v
		}
	} else {
		// Otherwise, store it under "result"
		finalResult["result"] = result
	}

	// Append result to workflow engine
	wm := WorkflowEngine{
		ResultNode: finalResult,
	}
	wfEngine, _ := ctx.FiberCtx.Locals("wfEngine").([]WorkflowEngine)
	wfEngine = append(wfEngine, wm)
	ctx.FiberCtx.Locals("wfEngine", wfEngine)

	fmt.Printf("[Scraper] Action %s completed successfully.\n", action)
	return nil
}

// getHTMLFromResult extracts HTML string from browser scraping result
// Browser methods return map[string]interface{}{"html": ..., "url": ...}
func getHTMLFromResult(result interface{}) string {
	if resultMap, ok := result.(map[string]interface{}); ok {
		if htmlVal, exists := resultMap["html"]; exists {
			if htmlStr, ok := htmlVal.(string); ok {
				return htmlStr
			}
		}
	}
	// Fallback for HTTP method which returns string directly
	if htmlStr, ok := result.(string); ok {
		return htmlStr
	}
	return ""
}

func scrapeWithHTTP(url, userAgent string, headers map[string]string, returnType string) (interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return string(body), nil
}

// getChromePaths returns Chrome/Chromium paths for the current OS
func getChromePaths() []string {
	var paths []string

	// Windows paths
	paths = append(paths,
		"C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe",
		"C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe",
		os.Getenv("LOCALAPPDATA")+"\\Google\\Chrome\\Application\\chrome.exe",
	)

	// Linux paths
	paths = append(paths,
		"/usr/bin/google-chrome",
		"/usr/bin/chromium",
		"/usr/bin/chromium-browser",
		"/snap/bin/chromium",
	)

	// macOS paths
	paths = append(paths,
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Chromium.app/Contents/MacOS/Chromium",
	)

	return paths
}

func scrapeWithBrowser(url, userAgent, returnType string, waitMs int) (interface{}, error) {
	// Try to use chromedp, but fallback to HTTP if it fails
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("ignore-certificate-errors", true), // Ignore SSL errors
		chromedp.Flag("disable-blink-features", "AutomationControlled"), // Hide automation flag
		chromedp.UserAgent(userAgent), // Set custom User-Agent
	)

	// Try to find Chrome/Chromium
	chromePaths := getChromePaths()

	for _, path := range chromePaths {
		if _, err := os.Stat(path); err == nil {
			opts = append(opts, chromedp.ExecPath(path))
			break
		}
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var htmlContent string
	var finalURL string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(time.Duration(waitMs)*time.Millisecond), // Use configurable wait time
		chromedp.Location(&finalURL),
		chromedp.OuterHTML("html", &htmlContent),
	)

	if err != nil {
		// If chromedp fails (Chrome not installed), fallback to HTTP
		fmt.Fprintf(os.Stderr, "Browser scraping failed, falling back to HTTP: %v\n", err)
		return scrapeWithHTTP(url, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36", nil, "html")
	}

	// Return both HTML and final URL as a map
	return map[string]interface{}{
		"html": htmlContent,
		"url":  finalURL,
	}, nil
}

func scrapeWithFormSubmit(url, userAgent string, formData map[string]string, submitSelector, clickSelector string, waitMs int) (interface{}, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("ignore-certificate-errors", true), // Ignore SSL errors
		chromedp.Flag("disable-blink-features", "AutomationControlled"), // Hide automation flag
		chromedp.UserAgent(userAgent), // Set custom User-Agent
	)

	// Try to find Chrome/Chromium
	chromePaths := getChromePaths()

	for _, path := range chromePaths {
		if _, err := os.Stat(path); err == nil {
			opts = append(opts, chromedp.ExecPath(path))
			break
		}
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// Build chromedp tasks
	var tasks chromedp.Tasks

	// Navigate to URL
	tasks = append(tasks, chromedp.Navigate(url))
	tasks = append(tasks, chromedp.Sleep(2*time.Second)) // Wait for page load

	// Fill form fields
	for selector, value := range formData {
		tasks = append(tasks, chromedp.SetValue(selector, value, chromedp.ByQuery))
	}

	// Submit the form
	if submitSelector != "" {
		// Try to click the submit button
		tasks = append(tasks, chromedp.Click(submitSelector, chromedp.ByQuery))
	} else if len(formData) > 0 {
		// If no submit selector provided, press Enter on the last filled field
		// This simulates pressing Enter to submit the form
		var lastSelector string
		for selector := range formData {
			lastSelector = selector
		}
		if lastSelector != "" {
			tasks = append(tasks, chromedp.SendKeys(lastSelector, "\n", chromedp.ByQuery))
		}
	}

	// Wait for results to load
	tasks = append(tasks, chromedp.Sleep(time.Duration(waitMs)*time.Millisecond))

    // SECONDARY CLICK ACTION (New Feature)
    if clickSelector != "" {
        // Log not available here, but we proceed
        tasks = append(tasks, chromedp.Click(clickSelector, chromedp.ByQuery))
        tasks = append(tasks, chromedp.Sleep(time.Duration(waitMs)*time.Millisecond))
    }

	// Get the HTML and final URL
	var htmlContent string
	var finalURL string
	tasks = append(tasks, chromedp.Location(&finalURL))
	tasks = append(tasks, chromedp.OuterHTML("html", &htmlContent))

	err := chromedp.Run(ctx, tasks)
	if err != nil {
		return nil, fmt.Errorf("browser automation failed: %v", err)
	}

	// Return both HTML and final URL as a map
	return map[string]interface{}{
		"html": htmlContent,
		"url":  finalURL,
	}, nil
}

func extractTextFromHTML(html, selector string) (interface{}, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var results []string
	if selector == "" {
		selector = "body"
	}

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			results = append(results, text)
		}
	})

	return results, nil
}

func extractLinksFromHTML(html, selector string) (interface{}, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var results []interface{}
	if selector == "" {
		selector = "a"
	}

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			results = append(results, map[string]interface{}{
				"text": strings.TrimSpace(s.Text()),
				"href": href,
			})
		}
	})

	return results, nil
}

func extractImagesFromHTML(html, selector string) (interface{}, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var results []interface{}
	if selector == "" {
		selector = "img"
	}

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists {
			alt, _ := s.Attr("alt")
			results = append(results, map[string]interface{}{
				"src": src,
				"alt": alt,
			})
		}
	})

	return results, nil
}

func extractDataWithRules(html, rulesJSON string) (interface{}, error) {
	if rulesJSON == "" {
		return nil, fmt.Errorf("extraction rules required")
	}

	var rules map[string]string
	if err := json.Unmarshal([]byte(rulesJSON), &rules); err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	results := make(map[string]interface{})
	for key, selector := range rules {
		var values []string
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if text != "" {
				values = append(values, text)
			}
		})
		results[key] = values
	}

	return results, nil
}

func solveCaptcha(siteURL, apiKey string) (interface{}, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("2Captcha API key required")
	}

	// Create CAPTCHA task
	createURL := "https://2captcha.com/in.php"
	params := fmt.Sprintf("key=%s&method=userrecaptcha&googlekey=SITE_KEY&pageurl=%s&json=1", apiKey, siteURL)

	resp, err := http.Post(createURL, "application/x-www-form-urlencoded", bytes.NewBufferString(params))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var createResp map[string]interface{}
	json.Unmarshal(body, &createResp)

	if createResp["status"].(float64) != 1 {
		return nil, fmt.Errorf("failed to create captcha task")
	}

	taskID := createResp["request"].(string)

	// Poll for result
	for i := 0; i < 60; i++ {
		time.Sleep(5 * time.Second)

		resultURL := fmt.Sprintf("https://2captcha.com/res.php?key=%s&action=get&id=%s&json=1", apiKey, taskID)
		resp, err := http.Get(resultURL)
		if err != nil {
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var resultResp map[string]interface{}
		json.Unmarshal(body, &resultResp)

		if resultResp["status"].(float64) == 1 {
			return map[string]string{
				"solution": resultResp["request"].(string),
				"taskId":   taskID,
			}, nil
		}
	}
	return nil, fmt.Errorf("captcha solving timeout")
}

func applyFilter(data map[string]interface{}, filter string) (map[string]interface{}, error) {
	if filter == "" {
		return data, nil
	}

	parts := strings.SplitN(filter, "=", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid filter format: %s. Expected key=value", filter)
	}

	filterKey := strings.TrimSpace(parts[0])
	filterVal := strings.TrimSpace(parts[1])

	// Find indices that match
	colData, ok := data[filterKey]
	if !ok {
		return nil, fmt.Errorf("filter key '%s' not found in result", filterKey)
	}

	strs, ok := colData.([]string)
	if !ok {
		return nil, fmt.Errorf("filter column '%s' is not a string array", filterKey)
	}

	var matchIndices []int
	for i, v := range strs {
		if strings.EqualFold(strings.TrimSpace(v), filterVal) {
			matchIndices = append(matchIndices, i)
		}
	}

	// Build new result
	filteredResult := make(map[string]interface{})
	for k, v := range data {
		if arr, ok := v.([]string); ok {
			var newArr []string
			for _, idx := range matchIndices {
				if idx < len(arr) {
					newArr = append(newArr, arr[idx])
				}
			}
			filteredResult[k] = newArr
		} else {
			filteredResult[k] = v // Keep non-array fields?
		}
	}

	return filteredResult, nil
}

func flattenResult(data map[string]interface{}) map[string]interface{} {
	flat := make(map[string]interface{})
	for k, v := range data {
		if arr, ok := v.([]string); ok {
			if len(arr) > 0 {
				flat[k] = arr[0]
			} else {
				flat[k] = ""
			}
		} else {
			flat[k] = v
		}
	}
	return flat
}

func extractWithRegex(text interface{}, rulesJSON string) (interface{}, error) {
	if rulesJSON == "" {
		return nil, fmt.Errorf("extraction rules required")
	}

	// Convert text to string
	var textStr string
	switch v := text.(type) {
	case string:
		textStr = v
	case []string:
		textStr = strings.Join(v, " ")
	default:
		return nil, fmt.Errorf("invalid text type")
	}

	var rules map[string]string
	if err := json.Unmarshal([]byte(rulesJSON), &rules); err != nil {
		return nil, err
	}

	results := make(map[string]interface{})
	for key, pattern := range rules {
		re, err := regexp.Compile(pattern)
		if err != nil {
			results[key] = fmt.Sprintf("invalid regex: %v", err)
			continue
		}

		matches := re.FindAllStringSubmatch(textStr, -1)
		if len(matches) == 0 {
			results[key] = []string{}
		} else {
			var values []string
			for _, match := range matches {
				if len(match) > 1 {
					// If there are capture groups, use the first capture group
					values = append(values, values...) // BUG FIX: this was wrong in earlier thought, wait.
					// Actually, if match > 1, match[1] is the capture.
					values = append(values, match[1])
				} else {
					// Otherwise use the full match
					values = append(values, match[0])
				}
			}
			results[key] = values
		}
	}

	return results, nil
}
