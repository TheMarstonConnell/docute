# Syntax

Every file in Docute is a simple markdown file ending with `.md`.

A sample layout of a docs page would live in the parent folder `docs`. For example, a subset of the Docute documentation filetree looks like this:
```
docs/
    syntax/
        syntax.md
    README.md
    SUMMARY.md
    logo.png
    colors.yaml
```

As you can see, everything lives in the `docs` folder, with this file residing in the `syntax` folder. To make these files accessible to the user, we add them to the `SUMMARY.md` file. The summary file outlines the navigation bar. An example of this summary file would look like this.

```markdown
# Docute

## Getting Started
* [Introduction](README.md)

## Syntax
* [Syntax](syntax/syntax.md)
```