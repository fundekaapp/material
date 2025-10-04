#!/usr/bin/env lua

-- Simple JSON decoder for extracting markdown field
local function extract_markdown(json_str)
    -- Find the "markdown" field and extract its value
    local markdown = json_str:match('"markdown"%s*:%s*"(.-)"[,}]')
    
    if not markdown then
        -- Try multiline pattern
        markdown = json_str:match('"markdown"%s*:%s*"(.+)"[,}]')
    end
    
    if not markdown then
        return nil, "Could not find 'markdown' field in JSON"
    end
    
    -- Unescape JSON string
    markdown = markdown:gsub('\\n', '\n')
    markdown = markdown:gsub('\\t', '\t')
    markdown = markdown:gsub('\\"', '"')
    markdown = markdown:gsub('\\\\', '\\')
    
    return markdown
end

local function file_exists(path)
    local f = io.open(path, "r")
    if f then
        f:close()
        return true
    end
    return false
end

local function read_file(path)
    local f = io.open(path, "r")
    if not f then
        return nil, "Could not open file: " .. path
    end
    local content = f:read("*all")
    f:close()
    return content
end

local function write_file(path, content)
    local f = io.open(path, "w")
    if not f then
        return false, "Could not write file: " .. path
    end
    f:write(content)
    f:close()
    return true
end

local function create_directory(path)
    os.execute("mkdir -p " .. path)
end

local function get_markdown_files(dir)
    local files = {}
    local p = io.popen('ls "' .. dir .. '"/*.md 2>/dev/null')
    if p then
        for file in p:lines() do
            table.insert(files, file)
        end
        p:close()
    end
    return files
end

local function get_basename(path)
    return path:match("([^/]+)$")
end

-- Main execution
local function main()
    local markdown_dir = "."
    local cleaned_dir = "cleaned"
    
    -- Check if we're in the markdown directory
    if not file_exists("../pdf") then
        print("Warning: Running from unexpected location. Expected to be in markdown/ folder.")
    end
    
    -- Create cleaned directory
    create_directory(cleaned_dir)
    print("Created/verified 'cleaned' directory")
    
    -- Get all markdown files
    local files = get_markdown_files(markdown_dir)
    
    if #files == 0 then
        print("No .md files found in current directory")
        return
    end
    
    print(string.format("Found %d markdown files to process\n", #files))
    
    local success_count = 0
    local error_count = 0
    
    -- Process each file
    for i, filepath in ipairs(files) do
        local basename = get_basename(filepath)
        
        -- Skip files in cleaned directory
        if not filepath:match("/cleaned/") then
            print(string.format("[%d/%d] Processing %s...", i, #files, basename))
            
            local content, err = read_file(filepath)
            if not content then
                print(string.format("  ✗ Error reading: %s", err))
                error_count = error_count + 1
            else
                -- Check if it's JSON (contains "markdown" field)
                if content:match('"markdown"%s*:') then
                    local markdown, extract_err = extract_markdown(content)
                    
                    if not markdown then
                        print(string.format("  ✗ Error: %s", extract_err))
                        error_count = error_count + 1
                    else
                        -- Write cleaned markdown
                        local output_path = cleaned_dir .. "/" .. basename
                        local write_ok, write_err = write_file(output_path, markdown)
                        
                        if write_ok then
                            print(string.format("  ✓ Cleaned and saved to %s", output_path))
                            success_count = success_count + 1
                        else
                            print(string.format("  ✗ Error writing: %s", write_err))
                            error_count = error_count + 1
                        end
                    end
                else
                    -- Already clean markdown, just copy
                    local output_path = cleaned_dir .. "/" .. basename
                    local write_ok, write_err = write_file(output_path, content)
                    
                    if write_ok then
                        print(string.format("  ✓ Already clean, copied to %s", output_path))
                        success_count = success_count + 1
                    else
                        print(string.format("  ✗ Error copying: %s", write_err))
                        error_count = error_count + 1
                    end
                end
            end
        end
    end
    
    print(string.format("\n✓ Successfully processed: %d", success_count))
    if error_count > 0 then
        print(string.format("✗ Errors: %d", error_count))
    end
    print("\nCleaned files are in: " .. cleaned_dir .. "/")
end

-- Run the script
main()