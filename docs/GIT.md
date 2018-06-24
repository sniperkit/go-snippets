# GIT

## Interesting
https://github.com/larsxschneider/git-repo-analysis

## References
- https://stackoverflow.com/questions/3313908/git-is-really-slow-for-100-000-objects-any-fixes
- https://mirrors.edge.kernel.org/pub/software/scm/git/docs/git-update-index.html
- https://stackoverflow.com/questions/4994772/ways-to-improve-git-status-performance/43644347#43644347
- https://mirrors.edge.kernel.org/pub/software/scm/git/docs/git-repack.html
- https://larsxschneider.github.io/2016/09/21/large-git-repos
- https://github.com/larsxschneider/git-repo-analysis
- https://vmiklos.hu/blog/sparse-checkout-example-in-git-1-7
- https://news.ycombinator.com/item?id=3548824
- https://blogs.msdn.microsoft.com/devops/2017/05/30/optimizing-git-beyond-gvfs/
- https://www.wired.com/2015/09/google-2-billion-lines-codeand-one-place/
- https://blogs.msdn.microsoft.com/devops/2017/05/30/optimizing-git-beyond-gvfs/
- https://docs.microsoft.com/en-us/azure/devops/git/git-at-scale
- https://www.youtube.com/watch?v=4XpnKHJAok8
- https://www.sitepoint.com/managing-huge-repositories-with-git/
- http://blog.prasoonshukla.com/mercurial-vs-git-scaling
- https://news.ycombinator.com/item?id=7648237
- https://medium.com/@maoberlehner/monorepos-in-the-wild-33c6eb246cb9
- https://developer.atlassian.com/blog/2015/10/monorepos-in-git/
- https://www.voucherify.io/blog/migration-to-a-monolithic-repository-without-the-code-freeze
- https://forums.resin.io/t/suggestions-for-deploying-from-a-monolithic-git-repo/467/8
- https://gist.github.com/arschles/5d7ba90495eb50fa04fc
- https://github.com/bluebird89/Docs/blob/0168a4f57b825a8565575780127663198ccbd3c6/Tools/Git.md
- https://confluence.atlassian.com/bitbucket/reduce-repository-size-321848262.html
- https://stackoverflow.com/questions/15752026/git-on-windows-shows-modified-files-all-the-time-even-for-newly-cloned-repo
- https://stackoverflow.com/questions/9750606/git-still-shows-files-as-modified-after-adding-to-gitignore
- https://stackoverflow.com/questions/14564946/git-status-shows-changed-files-but-git-diff-doesnt
- https://blog.andrewray.me/dealing-with-compiled-files-in-git/
- https://stackoverflow.com/questions/5787937/git-status-shows-files-as-changed-even-though-contents-are-the-same
- https://www.infoq.com/news/2017/02/GVFS
- http://www.codewrecks.com/blog/index.php/2017/06/22/optimize-your-local-git-repository-from-time-to-time/
- https://stackoverflow.com/questions/3119850/is-there-a-way-to-clean-up-git
- https://stackoverflow.com/questions/5613345/how-to-shrink-the-git-folder
- http://stevelorek.com/how-to-shrink-a-git-repository.html
- https://blog.github.com/2018-03-05-measuring-the-many-sizes-of-a-git-repository/
- https://git-scm.com/book/en/v2/Git-Internals-Maintenance-and-Data-Recovery
- http://erikimh.com/optimizing-a-previously-large-and-bloated-git-repository/
- https://www.atlassian.com/blog/git/handle-big-repositories-git
- https://beacots.com/optimize-size-of-git-repository/
- https://blogs.msdn.microsoft.com/bharry/2017/05/24/the-largest-git-repo-on-the-planet/
- https://www.infoq.com/news/2017/02/GVFS
- 

# Examples
3. Remove these files out of git(For ex: every .psd file in any folder): git filter-branch –index-filter ‘git rm –cached –ignore-unmatch **/*.psd’

4. Remove git backup: rm -rf .git/refs/original/

5. Expire all of loose objects: git reflog expire –expire=now –all

6. Check if ther’s any loose objects left: git fsck –full –unreachable

7. Repack(recompress objects file into .pack file): git repack -A -d

8. Run git garbage collector: git gc –aggressive –prune=now

9. Push changes: git push -f