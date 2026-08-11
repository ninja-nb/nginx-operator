# Wiki Sync Workflow

Use this flow to push Git-tracked pages from `wiki/` to GitHub Wiki.

## 1) Clone the wiki repository

```bash
git clone https://github.com/ninja-nb/nginx-operator.wiki.git /tmp/nginx-operator-wiki
```

## 2) Copy tracked pages from this repo

From the main repo root:

```bash
rsync -av --delete wiki/ /tmp/nginx-operator-wiki/
```

## 3) Commit and push to wiki

```bash
cd /tmp/nginx-operator-wiki
git add .
git commit -m "Update wiki ADR pages from repo mirror"
git push origin master
```

## 4) Verify

Open: `https://github.com/ninja-nb/nginx-operator/wiki`

## Notes

- GitHub wiki renders `Home.md` as the landing page
- Keep filenames stable to avoid broken wiki links
- Prefer editing pages under `wiki/` and syncing, rather than direct UI edits
