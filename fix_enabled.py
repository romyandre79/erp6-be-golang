import re

with open(r'c:\lara\www\erp6-be-golang\core\scheduler\scheduler.go', 'r', encoding='utf-8') as f:
    content = f.read()

# Find and replace the enabled check to be more flexible
old_check = 'if foundSchedule && enabled == "true" && cronExp != ""'
new_check = 'if foundSchedule && (enabled == "true" || enabled == "1" || strings.ToLower(enabled) == "true") && cronExp != ""'

content = content.replace(old_check, new_check)

with open(r'c:\lara\www\erp6-be-golang\core\scheduler\scheduler.go', 'w', encoding='utf-8') as f:
    f.write(content)

print("Fixed enabled check")
