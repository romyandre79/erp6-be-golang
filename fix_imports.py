with open(r'c:\lara\www\erp6-be-golang\core\scheduler\scheduler.go', 'r', encoding='utf-8') as f:
    lines = f.readlines()

# Find and fix the import section
new_lines = []
in_import = False
import_fixed = False

for i, line in enumerate(lines):
    if 'import (' in line:
        in_import = True
        new_lines.append(line)
        continue
    
    if in_import and not import_fixed:
        # Skip malformed lines and rebuild imports
        if ')' in line:
            # Insert correct imports before closing paren
            new_lines.append('\tgenerator "erp6-be-golang/core/generator/db"\n')
            new_lines.append('\t"erp6-be-golang/models"\n')
            new_lines.append('\t"encoding/json"\n')
            new_lines.append('\t"log"\n')
            new_lines.append('\t"strings"\n')
            new_lines.append('\t"sync"\n')
            new_lines.append('\t"time"\n')
            new_lines.append('\n')
            new_lines.append('\t"github.com/gofiber/fiber/v2"\n')
            new_lines.append('\t"github.com/robfig/cron/v3"\n')
            new_lines.append('\t"github.com/valyala/fasthttp"\n')
            new_lines.append('\t"gorm.io/gorm"\n')
            new_lines.append(line)
            in_import = False
            import_fixed = True
        continue
    
    new_lines.append(line)

with open(r'c:\lara\www\erp6-be-golang\core\scheduler\scheduler.go', 'w', encoding='utf-8') as f:
    f.writelines(new_lines)

print("Fixed imports")
