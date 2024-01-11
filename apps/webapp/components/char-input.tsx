import { useEffect, useRef, useState } from "react";
import './char-input.css';

export interface CharInputProps {
    // A string of the form 'xxx-xxx-xxx', where each 'x' will become an input
    template: string;

    onChange: (value: string) => void;
}

/**
 * CharInput is an input component inspired by 2FA screens, where each
 * character is its own input element, and typing a letter automatically 
 * moves you to the next input
 */
export function CharInput(props: CharInputProps) {
    const numValues = props.template.replace('-', '').length;
    const initialValue = Array.from(Array(numValues)).map(_ => ''); // one-liner to generate N sized array of empty strings

    const [values, setValues] = useState<string[]>(initialValue);
    useEffect(() => {
        const output = [];
        let i = 0;
        for(const char of props.template) {
            switch (char) {
                case 'X':
                case 'x':
                    output.push(values[i])
                    i++;
                    break;
                case '-':
                    output.push(['-'])
                    break;
            }
        }
        props.onChange(output.join(''));
    })

    const inputRefs = useRef<Array<HTMLInputElement | null>>([]);

    function focusIndex(prevIndex: number, nextIndex: number) {
        if (nextIndex >= 0 && nextIndex <= numValues - 1) {
            (inputRefs.current[nextIndex] as HTMLInputElement).focus()
        } else {
            (inputRefs.current[prevIndex] as HTMLInputElement).blur();
        }
    }

    function onKeyDown (i: number, e: KeyboardEvent) {
        e.preventDefault();

        if (e.key === 'Backspace') {
            const newValues = [...values]
            newValues[i] = '';
            setValues(newValues);

            focusIndex(i, i - 1);
        } else {
            const newValues = [...values]
            newValues[i] = e.key.toUpperCase();
            setValues(newValues);

            focusIndex(i, i + 1);
        }
    }

    const inputs = values.map((value, i) => (
        <input
            ref={ref => inputRefs.current[i] = ref}
            key={i}
            type="text"
            placeholder="X"
            value={value}
            onKeyDown={(e) => onKeyDown(i, e as any)}
            onChange={(e) => { }}
        ></input>
    ))

    const elements = []
    let inputCounter = 0;
    for (const char of props.template) {
        switch (char) {
            case 'X':
            case 'x':
                elements.push(inputs[inputCounter]);
                inputCounter++;
                break;
            case '-':
                elements.push(
                    <span key={"after-" + inputCounter}>-</span>
                )
                break;
            default:
                throw new Error("Template string can only contain 'X', 'x', and '-'.")
        }
    }

    return (
        <div className="char-input">
            {elements}
        </div>
    )
}