/*

Package entity defines database entity models for the insight/assistant ETL pipeline.

The entity types represent database tables at different stages of the data processing
pipeline, from raw source data through rule-based tagging to cleaned output.

# Entity Hierarchy

The entities form a layered hierarchy based on the stage of data processing:

  Raw → Rows → Data → DataContent → DataContentSegWords / DataContentSplitWords
                      ↓
                 Rule → DataRule → TagDataRule / TagDataRules
                      ↓
               CleanData / CleanRows

# Entity Types

## Core Data Entities

  Raw
    The most basic entity type. Maps directly to an existing database table without
    any naming transformation. Used as the entry point for the pipeline.

  Rows
    Intermediate row-level data table with a named primary key (idName). Supports
    cloning, suffix appending, and deleted-row tracking. Sits between Raw and Data
    in the pipeline.

  Data
    Processed data table whose actual database name is prefixed with "data_".
    Generated from Rows after row-level processing. Serves as input for rule-based
    analysis and tagging.

  DataContent
    A sub-table of a Data entity associated with a specific content field (e.g.,
    "body", "title"). Table name: "data_<data>_<content>".

## Word Analysis Entities

  DataContentSegWords
    Stores NLP-segmented words extracted from a content field (e.g., jieba).
    Table name: "tag_<data>_<content>_seg_words".

  DataContentSplitWords
    Stores delimiter-split words from a content field (e.g., whitespace, punctuation).
    Table name: "tag_<data>_<content>_split_words".

## Rule Entities

  Rule
    A rule definition consisting of a name, keyword configuration, and a set of labels.
    Used for keyword matching and label classification. Table name: "rule_<name>".
    Supports aliasing (same definition, different output names) and cloning.

  DataRule
    A Rule bound to a specific Data entity. Table name: "rule_<data>_<rule>".

## Tag Entities

  TagDataRule
    Tagging result for a single data+rule combination. Table name: "tag_<data>_<rule>".

  TagDataRules
    Tagging result for a data entity with multiple rules grouped under one name.
    Table name: "tag_<data>_<name>".

## Clean Entities

  CleanData
    Wraps the data cleaning lifecycle for a Data entity: source (Rows), cleaned
    output (Data), and deleted rows (Rows with "deleted_" prefix).

  CleanRows
    Wraps the data cleaning lifecycle for a Rows entity: source (Raw), rows being
    cleaned (Rows), and deleted rows (Rows with "deleted_" prefix).

# Table Naming Conventions

  Prefix/Suffix        Entity Type
  ─────────────────────────────────────
  (no prefix)          Raw, Rows
  data_                Data, DataContent
  rule_                Rule, DataRule
  tag_                 TagDataRule, TagDataRules,
                       DataContentSegWords, DataContentSplitWords
  deleted_             Deleted rows (derived from Rows)
  _seg_words           NLP-segmented words
  _split_words         Delimiter-split words

# Data Flow

  1. Raw/External Source
     Data enters the pipeline as a Raw entity, pointing to an existing database table.

  2. Row Processing
     Raw data is loaded into Rows entities for row-level operations. Rows can be
     cloned and suffixed for intermediate processing stages.

  3. Data Processing
     Rows are converted to Data entities ("data_*" tables) for further analysis.

  4. Content Analysis
     Data content fields can be further analyzed into:
       - DataContent: sub-tables for specific content fields
       - DataContentSegWords: NLP-segmented word frequency tables
       - DataContentSplitWords: delimiter-split word frequency tables

  5. Rule-Based Tagging
     Rules are defined (Rule) and optionally bound to data (DataRule), then applied
     to produce tag results (TagDataRule, TagDataRules with "tag_*" tables).

  6. Data Cleaning
     The pipeline supports cleaning stages via CleanData and CleanRows, which track
     both the cleaned output and the deleted rows separately.

# Architecture

The package uses a composition pattern with two embedded mixins:

  base
    Provides a shared database connection (*simpledb.SimpleDB) and a command hook
    mechanism (ExecCommand) for customizing schema operations before execution.

  extra
    Provides DDL/DML helper methods: AlterTable, Insert*, Drop, Truncate,
    CopyBuild, CopyBuildAndData, RawSqlAndScan, ExecSql.

All entity types implement one or more of the following interfaces:

  assistant.Raw       — basic table with Name() and TableName()
  assistant.Entity    — Raw with IDName()
  assistant.Rule      — rule definition with name, labels, keyword config
  assistant.RuleItems — rule items with aliases and keyword formatting
  etl.CleanResource   — clean stage with Source(), Data(), Deleted()
*/
package entity