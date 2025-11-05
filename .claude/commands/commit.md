---
description: Create a git commit with staged/unstaged changes after user confirmation
allowed-tools: Bash(git:*), Bash(task:*)
---

You are creating a git commit. Follow these steps carefully:

1. **Analyze the current git state:**
   - Run `git status` to see all untracked and modified files
   - Run `git diff` to see unstaged changes
   - Run `git diff --staged` to see staged changes
   - Run `git log --oneline -5` to see recent commit messages for context

2. **Draft a commit message:**
   - Analyze the changes and create a clear, concise commit message
   - Follow the conventional commit format if appropriate
   - **CRITICAL:** The first line MUST start with a lowercase letter (never capitalize the first word)
   - Keep the first line under 72 characters
   - Add a blank line, then bullet points explaining the changes
   - **CRITICAL:** Get the git user's name and email by running `git config user.name` and `git config user.email`
   - **CRITICAL:** Add a blank line at the end, then add: `Signed-off-by: [Name] <email>` using the git config values
   - **NEVER** include these lines in the commit message:
     - 🤖 Generated with [Claude Code](https://claude.com/claude-code)
     - Co-Authored-By: Claude <noreply@anthropic.com>

3. **Present the plan to the user:**
   Display in a clear format:

   ```
   ## Proposed Commit

   ### Files to be committed:
   [list staged files, or indicate which unstaged files will be staged]

   ### Commit Message:
   ```

   [show the full commit message including Signed-off-by line]

   ```

   ### Changed files preview:
   [show a summary of the key changes]
   ```

   Ask the user: **"Do you want to proceed with this commit? (yes/no)"**

4. **Wait for user confirmation:**
   - Do NOT proceed until the user explicitly responds with "yes", "y", or similar affirmation
   - If the user says "no" or wants changes, ask what they'd like to modify
   - If the user wants to edit the message, let them provide the new message

5. **Execute the commit ONLY after confirmation:**
   - **FIRST:** Run `task format` to format changed files before staging
   - If there are unstaged files that should be committed, run `git add [files]`
   - **Check if tests are needed:** Determine if any Go (.go) or TypeScript (.ts, .tsx) files are being committed
     - Run: `git diff --cached --name-only | grep -E '\.(go|ts|tsx)$'`
     - If the command returns any files, tests are needed
     - If no Go/TypeScript files are found, skip test execution and proceed to commit
   - **If tests are needed:**
     - **Run tests:** Execute `task test` to verify all tests pass
     - **If tests fail:**
       - Display the test failure output
       - Do NOT proceed with the commit
       - Ask the user: "Tests failed. What would you like to do? (fix/skip/cancel)"
       - If user says "skip", proceed with commit anyway (but warn them)
       - If user says "cancel", stop the commit process
       - If user says "fix", wait for them to fix and re-run the commit command
     - **If tests pass:** Proceed with the commit
   - **If tests not needed:** Proceed directly to commit
   - Run the git commit command with the confirmed message using a HEREDOC:
     ```bash
     git commit -m "$(cat <<'EOF'
     [commit message here with Signed-off-by line]
     EOF
     )"
     ```
   - If the commit is blocked by a pre-commit hook that formatted files, re-stage the formatted files and commit again
   - Run `git status` to confirm the commit succeeded
   - Display the commit hash and success message

**Important Rules:**

- NEVER commit without explicit user confirmation
- ALWAYS run tests before committing if Go or TypeScript files are modified
- ONLY commit if tests pass (unless user explicitly chooses to skip)
- ALWAYS include the Signed-off-by line with the user's git config information
- NEVER include Claude Code attribution or co-authorship lines
- If the commit fails, explain the error and ask how to proceed
