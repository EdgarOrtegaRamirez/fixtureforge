# Security Policy

## Overview
FixtureForge generates fake test data locally. It does not:
- Make network requests
- Store or transmit any data
- Execute arbitrary code
- Access files other than schema files you provide

## Data Privacy
All generated data is random and does not correspond to real individuals.
- Names are randomly generated from a fixed list
- Emails use random combinations with disposable-style domains
- Addresses, phone numbers, and other PII are synthetic

## Schema Files
Schema files are YAML/JSON configuration files processed locally.
They do not support code execution or arbitrary file access.

## Reporting
If you discover a security vulnerability, please open an issue on GitHub.
