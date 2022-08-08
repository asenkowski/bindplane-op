import {
  FormControl,
  Autocomplete,
  Chip,
  TextField,
  FormHelperText,
} from "@mui/material";
import { isArray, isFunction, isEmpty } from "lodash";
import { useState } from "react";
import { validateStringsField } from "../validation-functions";
import { useValidationContext } from "../ValidationContext";
import { ParamInputProps } from "./ParameterInput";

import styles from "./parameter-input.module.scss";

export const StringsParamInput: React.FC<ParamInputProps<string[]>> = ({
  definition,
  value,
  onValueChange,
}) => {
  const [inputValue, setInputValue] = useState("");
  const { setError, touched, errors, touch } = useValidationContext();

  // handleChipClick edits the selected chips value.
  function handleChipClick(ix: number) {
    if (!isArray(value)) {
      return;
    }

    // Edit the chips value
    setInputValue(value[ix]);

    // Remove the chip from the values because its being edited.
    const copy = [...value];
    copy.splice(ix, 1);
    isFunction(onValueChange) && onValueChange(copy);
  }

  // Make sure we "enter" the value if a user leaves the
  // input without hitting enter
  function handleBlur() {
    touch(definition.name);
    if (!isEmpty(inputValue)) {
      handleValueChange([...(value ?? []), inputValue]);
    }
  }

  function handleValueChange(newValue: string[]) {
    onValueChange && onValueChange(newValue);

    setInputValue("");
    setError(
      definition.name,
      validateStringsField(newValue, definition.required)
    );
  }

  const label = definition.required
    ? `${definition.label} *`
    : `${definition.label}`;

  return (
    <FormControl fullWidth>
      <Autocomplete
        options={[]}
        multiple
        disableClearable
        freeSolo
        classes={{
          root: definition.relevantIf ? styles.indent : undefined,
        }}
        // value and onChange pertain to the string[] value of the input
        value={value ?? []}
        onChange={(e, v: string[]) => handleValueChange(v)}
        // inputValue and onInputChange refer to the latest string value being entered
        inputValue={inputValue}
        onInputChange={(e, newValue) => setInputValue(newValue)}
        onBlur={handleBlur}
        renderTags={(value: readonly string[], getTagProps) =>
          value.map((option: string, index: number) => (
            <Chip
              size="small"
              variant="outlined"
              label={option}
              {...getTagProps({ index })}
              classes={{ label: styles.chip }}
              onClick={() => handleChipClick(index)}
            />
          ))
        }
        renderInput={(params) => (
          <TextField
            {...params}
            label={label}
            size={"small"}
            helperText={definition.description}
            id={definition.name}
          />
        )}
      />
      {touched[definition.name] && errors[definition.name] && (
        <FormHelperText error={true}>{errors[definition.name]}</FormHelperText>
      )}
    </FormControl>
  );
};
