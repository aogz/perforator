# perforator

Developer statistics tracker

![image](https://user-images.githubusercontent.com/10484630/139407418-dd56e2f1-377d-4b55-9863-db1d1012d0b5.png)

## Available stats / commands:
- `rejection-rate`: Percentage of rejected pull. Shows aggregated values per repository and value per developer

## Usage:
1. Add `GITHUB_ACCESS_TOKEN` env var
2. Run cli tool:
```bash
perforator rejection-rate --repo django/django --limit 100
```
