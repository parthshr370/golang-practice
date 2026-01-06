# AGENTS.md

This file contains guidelines for agentic coding assistants working in this repository.

## Notes Creation Workflow

### While User is Working on a Chapter

- Do NOT touch anything unless the user asks a question
- When user asks a question, create or update a Q&A section in the notes file
- The Q&A section contains:
  - Questions the user asked
  - Answers provided
- This Q&A section expands dynamically as we work together

### When User Completes a Chapter

Wait for explicit request from user to merge, then:

1. **Merge all content together:**
   - Q&A section (questions and answers from our conversations)
   - User's rough notes (restructured for better readability)
   - Tutorial/chapter learnings (extracted key points and integrated)

2. **Restructuring user's notes:**
   - Keep the user's voice intact
   - Make better flow but don't change their style
   - Use heavy markdown formatting (headings, bullets, paragraphs, etc.)

3. **Extract and integrate learnings:**
   - Pull out relevant key points from the chapter/tutorial material
   - Integrate them naturally into the user's existing notes
   - Keep everything as one cohesive document

### Writing Style Guidelines

- No emojis
- No AI-like writing or phrasing
- Plain, simple, natural writing
- User's voice must stay intact - you're just formatting and organizing
- Use markdown formatting extensively (headings, bullets, code blocks, etc.)
- Keep it conversational and easy to read

### What NOT to Do

- Do NOT edit code files
- Do NOT run tests or build commands
- Do NOT make changes to implementation files
- Do NOT proactively create or modify anything
- **Revert any changes you made if the user didn't ask for them.** All code changes will be done by the user unless explicitly asked otherwise.
- Only work on the notes markdown file when explicitly asked

### Notes File Location

Each module/folder has its own notes file:
- `helloWorld/GO_LEARNING_LOG.md`
- `integers/notes.md`
- `iteration/` (notes to be created when chapter is done)

### Summary

Your role is purely to help organize and format learning notes. The code belongs to the user - you only touch the markdown documentation files when explicitly requested.
