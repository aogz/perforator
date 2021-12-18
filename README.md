# Perforator

Perforator is a CLI tool to track and analyze performance and statistics based on the information from GitHub Pull Requests. 

Data received using this tool should be taken with a grain of salt, mostly because every single PR is different and lots of factors affect PR's review time, reasons for rejecting changes, etc.

![image](https://user-images.githubusercontent.com/10484630/139407418-dd56e2f1-377d-4b55-9863-db1d1012d0b5.png)

## Prerequisites
This tool heavily relies on Github API, so [GITHUB_ACCESS_TOKEN](https://github.com/settings/tokens) environment variable is required to make requests to the API. Depending on the command, different permissions might be needed, but normally `repo` permissions should be enough.

## Available stats / commands:
- `rejection-rate`: Percentage of rejected pull. Shows aggregated values per repository and value per developer

    Options:
    - `repo`: repository name in `owner/name` format, e.g. `facebook/react` or `django/django`.
    - `limit` (default: 10): number of PRs/tickets to collect data from.
    - `skip` (default: 0): number of PRs/tickets to skip from the beggining of the list.
    - `contributors`: comma-separated list of contributors, e.g. aogz,foo,bar

    Examples:
    ```bash
    $ perforator rejection-rate --repo django/django --limit 100
    ```

- `review-time`: Command to track time PR was in review.

    Options:
    - `repo`: repository name in `owner/name` format, e.g. `facebook/react` or `django/django`.
    - `limit` (default: 10): number of PRs/tickets to collect data from.
    - `skip` (default: 0): number of PRs/tickets to skip from the beggining of the list.
    - `contributors`: comma-separated list of contributors, e.g. aogz,foo,bar
    - `group-by` (default: `reviewer`): criteria to group by. Can be one of: `author`, `reviewer`:
        
        `author`: Calculates average time PRs from specific author spends from the moment pr was created to the moment it was merged.
        
        `reviewer`: Calculates average time for a specific reviewer to review PRs.

    Examples:
    ```bash
    $ perforator review-time --repo facebook/react --limit 10 --group-by author 
    ```

    ```bash
    $ perforator review-time --repo aogz/perforator --limit 10 --group-by reviewer 
    ```

- `issue-author`: Number of issues created by author

    Options:
    - `repo`: repository name in `owner/name` format, e.g. `facebook/react` or `django/django`.
    - `limit` (default: 10): number of PRs/tickets to collect data from.
    - `skip` (default: 0): number of PRs/tickets to skip from the beggining of the list.
    - `contributors`: comma-separated list of contributors, e.g. aogz,foo,bar

    Examples:
    ```bash
    $ perforator issue-author --repo django/django --limit 100
    ```

- `issue-labels`: Number of issues grouped by labels

    Options:
    - `repo`: repository name in `owner/name` format, e.g. `facebook/react` or `django/django`.
    - `limit` (default: 10): number of PRs/tickets to collect data from.
    - `skip` (default: 0): number of PRs/tickets to skip from the beggining of the list.
    - `contributors`: comma-separated list of contributors, e.g. aogz,foo,bar

    Examples:
    ```bash
    $ perforator issue-labels --repo django/django --limit 100
    ```
