---
title: "Hugo Throwback"
date: 2018-03-11T02:33:20+08:00
tags: ["experiment", "go"]
draft: true
---

1. using hugo-theme-cactus-plus
    - enable tags by:
        1. adding `tag` in `taxonomies` in config file, ie.

        ```
        [taxonomies]
            tag = "tags"
        ``` 
        2. add `tags` in page's frontmatter. ie. `tags = ["go"]`
        3. visit posts by tags from `/tags/go`
2. frontmatter
    - use `=` with `+++` divider (toml), or `:` when `---` is used (yaml)
3. still build when there's error (no error display)
4. need to read docs thoroughly to understand the features. the docs is awesome tho!
5. to add translation:
    - add `language` and its weight in config, ie.

    ```
    [languages.en]
    weight = 1

    [languages.id]
        weight = 2
    ```
    - create post with language suffix, ie. `about/_index.id.md`