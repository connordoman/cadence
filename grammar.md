# Grammar

## Syntactic Grammar

```ebnf
expression
  -> "EVERY" ( implicit_schedule | explicit_schedule ) [ date_range ] ;

date_range
  -> "FROM" DATE [ inclusivity ] [ "TO" DATE [ inclusivity ]] ;

inclusivity
  -> "INCLUSIVE" | "EXCLUSIVE" | "INC" | "EXC"

implicit_schedule
  -> selector ;

explicit_schedule
  -> interval_spec [ "ON" selector ] ;

interval_spec
  -> [ INTEGER ] unit ;

unit
  -> "DAY" | "DAYS"
   | "WEEK" | "WEEKS"
   | "MONTH" | "MONTHS" ;

selector
  -> "WEEKDAYS"
   | day_list
   | ordinal_day_list ;

day_list
  -> day ( "," day )* ;

ordinal_day_list
  -> ordinal_day ( "," ordinal_day )* ;

ordinal_day
  -> distinction day ;

distinction
  -> "FIRST" | "LAST" ;

day
  -> "MON" | "TUE" | "WED" | "THU" | "FRI" | "SAT" | "SUN" ;
```
