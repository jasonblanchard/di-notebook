resource "aws_glue_catalog_database" "di_notebook_production" {
  name = "di_notebook_production"
}

resource "aws_glue_catalog_table" "entry_revisions" {
  name       = "di_entry_revisions"
  database_name = aws_glue_catalog_database.di_notebook_production.name
  table_type = "EXTERNAL_TABLE"

  storage_descriptor {
    location      = "s3://di-entry-revisions-production/"
    input_format  = "org.apache.hadoop.hive.ql.io.parquet.MapredParquetInputFormat"
    output_format = "org.apache.hadoop.hive.ql.io.parquet.MapredParquetOutputFormat"

    ser_de_info {
      serialization_library = "org.apache.hadoop.hive.ql.io.parquet.serde.ParquetHiveSerDe"
      parameters = {
        "serialization.format" = 1
      }
    }

    columns {
      name = "id"
      type = "string"
    }

    columns {
      name = "text"
      type = "string"
    }

    columns {
      name = "creator_id"
      type = "string"
    }

    columns {
      name = "create_time"
      type = "timestamp"
    }

    columns {
      name = "update_time"
      type = "timestamp"
    }

    columns {
      name = "delete_time"
      type = "timestamp"
    }

    columns {
      name = "actor_type"
      type = "string"
    }

    columns {
      name = "actor_id"
      type = "string"
    }
  }

  partition_keys {
    name = "year"
    type = "string"
  }

  partition_keys {
    name = "month"
    type = "string"
  }
}