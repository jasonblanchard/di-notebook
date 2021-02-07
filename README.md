# di-notebook

Holds all the notebook logic

Order of operations to deprecate is_deleted
1. ~Migration to add delete_time column~
2. ~Deploy ^~
3. ~Code to dual write delete_time & is_deleted~
4. ~Deploy ^~
5. Script data migration
  - if is_deleted == true AND delete_time is nil
  - set delete_time to updated_at
6. Remove dual write code & migration to remove is_deleted column