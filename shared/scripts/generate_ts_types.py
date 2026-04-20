import os
import subprocess
import re
import glob
import argparse

def parse_args():
    parser = argparse.ArgumentParser(description="Generate TS models and enums from JSON schemas.")
    parser.add_argument("--schema-dir", default="../schemas", help="Directory containing JSON schemas")
    parser.add_argument("--enums-dir", default="./types/enums", help="Output directory for enums")
    parser.add_argument("--models-dir", default="./types/models", help="Output directory for models")
    return parser.parse_args()

def camel_to_kebab(name):
    name = re.sub('(.)([A-Z][a-z]+)', r'\1-\2', name)
    return re.sub('([a-z0-9])([A-Z])', r'\1-\2', name).lower()

def ensure_dirs(enums_dir, models_dir):
    for d in [enums_dir, models_dir]:
        if os.path.exists(d):
            for f in glob.glob(os.path.join(d, "*.ts")): os.remove(f)
        else:
            os.makedirs(d)

def format_ts_directories(enums_dir, models_dir):
    print("🎨 Formatting files with Prettier...")
    try:
        subprocess.run(["npx", "--yes", "prettier", "--write", f"{enums_dir}/*.ts", f"{models_dir}/*.ts"], 
                       check=True, capture_output=True)
    except Exception: pass

def run_quicktype(tmp_file, schema_dir):
    schema_files = glob.glob(os.path.join(schema_dir, "*.json"))
    if not schema_files: return False
    cmd = ["quicktype", "--src-lang", "schema", "--lang", "ts", "--just-types", "-o", tmp_file]
    cmd.extend(schema_files)
    subprocess.run(cmd)
    return True

def extract_interfaces(ts_code):
    interfaces = {}
    for match in re.finditer(r'export interface (\w+) \{', ts_code):
        name = match.group(1)
        start_idx = match.end() - 1 
        brace_count = 0
        end_idx = -1
        for i in range(start_idx, len(ts_code)):
            if ts_code[i] == '{': brace_count += 1
            elif ts_code[i] == '}':
                brace_count -= 1
                if brace_count == 0:
                    end_idx = i
                    break
        if end_idx != -1:
            interfaces[name] = ts_code[start_idx+1 : end_idx]
    return interfaces

def get_clean_name(raw_name):
    """Transforms StoryboardEntityObject -> IStoryboardEntity"""
    clean = raw_name
    # Remove suffixes
    for suffix in ["Object", "Element", "Class"]:
        if clean.endswith(suffix) and len(clean) > len(suffix):
            clean = clean[:-len(suffix)]
    # Add I prefix
    if not (clean.startswith("I") and any(c.isupper() for c in clean[1:2])):
        clean = f"I{clean}"
    return clean

def process_and_split(tmp_file, args):
    with open(tmp_file, "r") as f:
        full_content = f.read()

    # 1. DISCOVER ENUMS
    enum_matches = list(re.finditer(r'export enum (\w+) \{([^}]+)\}', full_content))
    discovered_enums = {m.group(1): re.findall(r'"([^"]+)"', m.group(2)) for m in enum_matches}
    
    # 2. DISCOVER STRUCTS & BUILD GLOBAL RENAME MAP
    raw_structs = extract_interfaces(full_content)
    rename_map = {name: get_clean_name(name) for name in raw_structs.keys()}
    
    # Identify duplicates (e.g., StoryboardEntityObject vs StoryboardEntity)
    # If two raw names map to the same clean name, we keep only one
    unique_structs = {}
    for raw_name, clean_name in rename_map.items():
        if clean_name not in unique_structs:
            unique_structs[clean_name] = raw_structs[raw_name]
        else:
            # If we find a shorter raw name, it's usually the 'primary' one
            if len(raw_name) < len([k for k, v in rename_map.items() if v == clean_name][0]):
                unique_structs[clean_name] = raw_structs[raw_name]

    print(f"🔍 TS Discovered: {len(discovered_enums)} Enums, {len(unique_structs)} Clean Interfaces")

    # 3. GENERATE ENUMS
    for name, values in discovered_enums.items():
        filename = f"{camel_to_kebab(name)}.enum.ts"
        union = "\n".join([f'  | "{v}"' for v in values])
        with open(os.path.join(args.enums_dir, filename), "w") as f:
            f.write(f"export type {name} =\n{union};\n")

    # 4. GENERATE MODELS
    for clean_name, body in unique_structs.items():
        if clean_name in ["IType", "IModels", "IEnums", "IBaseItem"]: continue

        lines = body.split('\n')
        cleaned_lines = []
        is_entity = False
        used_enums = set()
        used_models = set()

        for line in lines:
            trimmed = line.strip()
            if not trimmed or "[property: string]: any;" in trimmed: continue

            # STRIP DYNAMO KEYS
            if re.search(r'\b(Pk|Sk|Gsi1Pk|Gsi1Sk)\b\s*\??\s*:', trimmed):
                is_entity = True
                if cleaned_lines and cleaned_lines[-1].strip() == "*/":
                    cleaned_lines.pop()
                    while cleaned_lines and cleaned_lines[-1].strip() != "/**": cleaned_lines.pop()
                    if cleaned_lines: cleaned_lines.pop()
                elif cleaned_lines and cleaned_lines[-1].strip().startswith("//"):
                    cleaned_lines.pop()
                continue

            # RENAME TYPES IN LINE BODY
            # This replaces 'StoryboardEntityObject' with 'IStoryboardEntity' everywhere
            for raw_m, clean_m in rename_map.items():
                if re.search(rf"\b{raw_m}\b", line):
                    line = re.sub(rf"\b{raw_m}\b", clean_m, line)
                    if clean_m != clean_name: used_models.add(clean_m)

            for enum_n in discovered_enums.keys():
                if re.search(rf"\b{enum_n}\b", line): used_enums.add(enum_n)

            cleaned_lines.append(line)

        if is_entity: cleaned_lines.insert(0, "  id: string;")
        
        # Build Imports
        imports = []
        for e in sorted(used_enums):
            imports.append(f'import type {{ {e} }} from "../enums/{camel_to_kebab(e)}.enum";')
        for m in sorted(used_models):
            # Convert clean name IName back to kebab name for file path
            kebab_m = camel_to_kebab(m[1:] if m.startswith("I") else m)
            imports.append(f'import type {{ {m} }} from "./{kebab_m}.model";')

        file_content = "\n".join(sorted(list(set(imports)))) + "\n\n" if imports else ""
        file_content += f"export interface {clean_name} {{\n" + "\n".join(cleaned_lines) + "\n}\n"

        # Save file based on clean name
        base_name = clean_name[1:] if clean_name.startswith("I") else clean_name
        with open(os.path.join(args.models_dir, f"{camel_to_kebab(base_name)}.model.ts"), "w") as f:
            f.write(file_content)

if __name__ == "__main__":
    args = parse_args()
    ensure_dirs(args.enums_dir, args.models_dir)
    TEMP = "raw_ts.tmp"
    if run_quicktype(TEMP, args.schema_dir):
        process_and_split(TEMP, args)
        if os.path.exists(TEMP): os.remove(TEMP)
        format_ts_directories(args.enums_dir, args.models_dir)
        print("\n✨ Generation Complete.")