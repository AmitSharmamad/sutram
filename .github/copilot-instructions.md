# Clean Code Architecture Instructions

When writing code, adhere to the principles of Clean Code Architecture while maintaining the existing folder and module structure. Follow these guidelines:

1. **Separation of Concerns**:
    - Ensure each module has a single responsibility.
    - Avoid mixing business logic with infrastructure or UI code.

2. **Dependency Rule**:
    - Inner layers (core logic) should not depend on outer layers (UI, database, etc.).
    - Use interfaces or abstractions to invert dependencies.

3. **Modularity**:
    - Keep modules self-contained and reusable.
    - Avoid circular dependencies between modules.

4. **Readability**:
    - Write clear and concise code with meaningful names for variables, functions, and classes.
    - Add comments only where necessary to explain complex logic.

5. **Testing**:
    - Write unit tests for core logic and integration tests for module interactions.
    - Ensure tests do not depend on external systems unless explicitly required.

6. **Folder Structure**:
    - Respect the existing folder structure.
    - Place new files in appropriate directories based on their purpose.

7. **Code Reviews**:
    - Ensure all code adheres to these principles during code reviews.
    - Refactor legacy code incrementally to align with Clean Code Architecture.

By following these instructions, you can write maintainable and scalable code without disrupting the current structure.

Do not write tests unless specicially mentioned in the prompt.

Do not hallucinate, if you have any doubts or things are not clear, ask me.