---
name: README updater
description: Update README.md with the current project tree, commit the documentation change, and push it to the configured Git remote.
---

# README Updater

Update `README.md` to reflect the current repository contents whenever this agent is run.

## Workflow

1. Inspect the workspace recursively before editing. Include source files, configuration files, data files, and important project metadata.
2. Exclude `.git/`, editor or OS metadata, dependency caches, build output, temporary files, and other generated directories. Include `.github/` configuration files when they are relevant to maintaining the project.
3. Read each included text file enough to write an accurate description in one or two sentences. Do not guess from filenames alone.
4. Preserve any useful human-written introduction or sections already present in `README.md`.
5. Replace the content between the markers `<!-- BEGIN PROJECT STRUCTURE -->` and `<!-- END PROJECT STRUCTURE -->`. If the markers do not exist, add the generated section after the existing introduction.
6. Use this structure:

   ```markdown
   <!-- BEGIN PROJECT STRUCTURE -->
   ## Project Structure

   ```text
   .
   |-- directory/
   |   `-- file.ext
   `-- file.ext
   ```

   ### `path/to/file.ext`
   One or two concise sentences describing the file's purpose and important behavior.
   <!-- END PROJECT STRUCTURE -->
   ```

7. List directories in tree order and list files under their containing directory. Keep the tree readable and do not include the README's generated description for itself.
8. Keep descriptions factual, concise, and current. Mention important public routes, data formats, startup behavior, or generated/runtime status when applicable.
9. Do not change application source code, dependency files, generated databases, or unrelated documentation. Only edit `README.md` unless a README update genuinely requires a small documentation-only adjustment.
10. Preserve the existing README encoding when possible. If the editor cannot preserve it, use UTF-8 consistently and keep the Markdown valid.

## Commit and Push

After the README update and verification:

1. Check whether `README.md` changed. If it did not change, do not create a commit or push anything.
2. Check that the repository has a configured `origin` remote and that the current branch has a name. If either is missing, leave the README update in place and report the prerequisite instead of committing or pushing.
3. Stage only `README.md`; never use `git add .` or `git add -A`.
4. Create a commit with the exact message `docs: update README project structure`.
5. Push the current branch to `origin` with `git push origin HEAD`.
6. Never force-push, amend an existing commit, rewrite history, or include unrelated staged or unstaged files.
7. If commit or push fails because of authentication, branch protection, rejected remote changes, or another Git error, keep the local README change and report the error without retrying destructively.

## Verification

- Re-read the final `README.md`.
- Confirm every included project file appears exactly once in the tree and has a matching description.
- Confirm excluded generated or environment-specific files are absent.
- For this Go repository, run `go test ./...` when available and report the result.