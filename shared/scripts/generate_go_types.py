import os
import subprocess
import re
import glob
import argparse

def parse_args():
    parser = argparse.ArgumentParser(description="Generate Go models and enums from JSON schemas.")
    parser.add_argument("--schema-dir", default="../schemas", help="Directory containing JSON schemas")
    parser.add_argument("--enums-dir", default="./enums", help="Output directory for enums")
    parser.add_argument("--models-dir", default="./models", help="Output directory for models")
    parser.add_argument("--enums-import", default="github.com/gianghp123/Vidmerce/backend/internal/core/enums", help="Go import path for enums")
    return parser.parse_args()

def camel_to_snake(name):
    name = re.sub('(.)([A-Z][a-z]+)', r'\1_\2', name)
    return re.sub('([a-z0-9])([A-Z])', r'\1_\2', name).lower()

def to_pascal_case(s):
    return "".join(word.capitalize() for word in s.lower().replace("-", "_").split("_"))

def ensure_dirs(enums_dir, models_dir):
    for d in [enums_dir, models_dir]:
        if os.path.exists(d):
            for f in glob.glob(os.path.join(d, "*.go")): os.remove(f)
        else:
            os.makedirs(d)

def format_go_file(filepath):
    try:
        subprocess.run(["go", "fmt", filepath], check=True, capture_output=True)
    except Exception: pass

def run_quicktype(tmp_file, schema_dir):
    schema_files = glob.glob(os.path.join(schema_dir, "*.schema.json"))
    if not schema_files: 
        schema_files = glob.glob(os.path.join(schema_dir, "*.json"))
    if not schema_files: return False
    cmd = ["quicktype", "--src-lang", "schema", "--lang", "go", "--just-types", "--field-tags", "dynamodbav", "-o", tmp_file]
    cmd.extend(schema_files)
    subprocess.run(cmd)
    return True

def extract_structs(go_code):
    structs = {}
    for match in re.finditer(r'type (\w+) struct \{', go_code):
        name = match.group(1)
        start_idx = match.end() - 1
        brace_count = 0
        end_idx = -1
        for i in range(start_idx, len(go_code)):
            if go_code[i] == '{': brace_count += 1
            elif go_code[i] == '}':
                brace_count -= 1
                if brace_count == 0:
                    end_idx = i
                    break
        if end_idx != -1:
            structs[name] = go_code[start_idx+1 : end_idx]
    return structs

def get_clean_struct_name(raw_name):
    clean = raw_name
    for suffix in ["Object", "Element", "Class"]:
        if clean.endswith(suffix) and len(clean) > len(suffix):
            clean = clean[:-len(suffix)]
    if clean == "BaseItem": return "BaseItem"
    if not clean.endswith("Entity"):
        clean = f"{clean}Entity"
    return clean

def process_and_split(tmp_file, args):
    with open(tmp_file, "r") as f:
        full_content = f.read()

    enum_matches = list(re.finditer(r'type (\w+) string\s+const \((.*?)\)', full_content, re.DOTALL))
    discovered_enums = {m.group(1): m.group(2) for m in enum_matches}
    
    raw_structs = extract_structs(full_content)
    rename_map = {name: get_clean_struct_name(name) for name in raw_structs.keys()}
    
    unique_structs = {}
    for raw_name, clean_name in rename_map.items():
        if clean_name not in unique_structs:
            unique_structs[clean_name] = raw_structs[raw_name]
        else:
            if len(raw_name) < len([k for k, v in rename_map.items() if v == clean_name][0]):
                unique_structs[clean_name] = raw_structs[raw_name]

    for name, body in discovered_enums.items():
        filename = f"{camel_to_snake(name)}.enum.go"
        filepath = os.path.join(args.enums_dir, filename)
        fixed_lines = []
        for line in body.split('\n'):
            match = re.search(rf'(\w+)\s+{name}\s+=\s+"([^"]+)"', line)
            if match:
                val = match.group(2)
                new_const_name = f"{name}{to_pascal_case(val)}"
                fixed_lines.append(f'\t{new_const_name} {name} = "{val}"')
            elif line.strip(): fixed_lines.append(line)
        with open(filepath, "w") as f:
            f.write(f"package enums\n\ntype {name} string\n\nconst (\n" + "\n".join(fixed_lines) + "\n)\n")
        format_go_file(filepath)

    for struct_name, body in unique_structs.items():
        if struct_name in ["Type", "Models", "Enums", "TypeEntity", "ModelsEntity"]: 
            continue

        # --- UPDATED CASE-INSENSITIVE ENTITY DETECTION ---
        is_entity = re.search(r'dynamodbav:"(pk|sk|gsi1pk|gsi1sk)"', body, re.IGNORECASE) is not None and struct_name != "BaseItem"

        lines = body.split('\n')
        cleaned_lines = []
        uses_enum = False

        for line in lines:
            trimmed = line.strip()
            if not trimmed: continue
            
            if ' int64 ' in line:
                line = line.replace(' int64 ', ' int ')


            # Force non-pointers for system keys
            if re.search(r'\b(Gsi1Pk|Gsi1Sk|Pk|Sk)\b\s+\*string', line, re.IGNORECASE):
                line = re.sub(r'(\b(?:Gsi1Pk|Gsi1Sk|Pk|Sk)\b)\s+\*string', r'\1 string', line, flags=re.IGNORECASE)

            # --- UPDATED CASE-INSENSITIVE STRIPPING ---
            if is_entity:
                if re.search(r'dynamodbav:"(pk|sk|gsi1pk|gsi1sk)(,omitempty)?"', line, re.IGNORECASE):
                    if cleaned_lines and cleaned_lines[-1].strip().startswith("//"):
                        cleaned_lines.pop()
                    continue

            for raw_m, clean_m in rename_map.items():
                if re.search(rf"\b{raw_m}\b", line):
                    line = re.sub(rf"\b{raw_m}\b", clean_m, line)

            for enum_n in discovered_enums.keys():
                pattern = rf"(\s+(?:\[\])?\*?\b){enum_n}(\s+`dynamodbav)"
                if re.search(pattern, line):
                    line = re.sub(pattern, rf"\1enums.{enum_n}\2", line)
                    uses_enum = True
            
            cleaned_lines.append(line)

        final_body = "\n".join(cleaned_lines)
        if is_entity:
            final_body = "\n\tBaseItem\n" + final_body

        name_for_file = struct_name[:-6] if struct_name.endswith("Entity") and struct_name != "Entity" else struct_name
        filename = f"{camel_to_snake(name_for_file)}.model.go"
        filepath = os.path.join(args.models_dir, filename)
        
        import_stmt = f'\nimport "{args.enums_import}"\n' if uses_enum else ""
        file_content = f"package models\n{import_stmt}\ntype {struct_name} struct {{{final_body}\n}}\n"
        
        with open(filepath, "w") as f: f.write(file_content)
        format_go_file(filepath)
        
    print(f"✅ Generation complete.")

if __name__ == "__main__":
    args = parse_args()
    ensure_dirs(args.enums_dir, args.models_dir)
    TEMP = "raw_go.tmp"
    if run_quicktype(TEMP, args.schema_dir):
        process_and_split(TEMP, args)
        if os.path.exists(TEMP): os.remove(TEMP)