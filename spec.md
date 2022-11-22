# ARGON 3 dev 1.0

This is the specification on how to create an Argon interpreter.

This specification should be used as a guildline, and is subject to change for later versions.

Later updates for Argon 3 should be backwards compatable (where possible) with code designed for older versions
of the interpreter.

The Argon 3 interpreter is not intentionally designed to understand code written for Argon 1 and/or 2, but may
follow some design patterns found in those interpreters.

Argon 3 is a programming language primarily designed for maths computation, it is designed to use techniques and
rules set in maths. It's designed to be easy for mathematicians to write and understand algorithms in.
(e.g. `f(x) = x^2` to set a function)

reused variables, and infomation for use in understanding the pseudo REGEX:

-   NAME = [a-zA-Z][a-za-z0-9]\*
-   spaces used in the pseudo REGEX should be taken as 1 or many spaces.
-   there can be 0 or more spaces around any of the pseudo REGEX

---

## set variable

`(let/const or nothing) {NAME} = {ANY}`

let and const variables will set variables. if a const is already set in that stack, it will throw an error.
at the end of a opperation (e.g. if, while, or function) drop the variable stack.

having no verb at the start suggests the program would like to edit a variables from a different stack.
if there is no variable found in the other stacks, then it sets one in the current stack as a let.

setting variables returns the value, which can be used.

example:

```
if (x = 10) > 5 [
    log(x, 'is bigger than 5')
]
```

---

## functions

`(let/const or nothing) {NAME}({PARAMS}) = {CODE}`

the starting verb follows the rules set by [set variable](#set-variable).

function can be called by using its name followed by brackets with the params.

example:

```

const f(x) = x^2 + 2*x + 1
log('f(10) =', f(10))

```

output:

```

f(10) = 121

```

if the function does not return, then the value given is unknown/null

---

## wrap

`[{CODE}]`

a wrap encloses code in square brackets. its used to create a new scope, so variables set from
inside the wraps scope are deleted once the wrap is finished.

example:

```
let name = unknown
[
    name = input('name: ')
    let age = input('age: ')
    log('you are', age, 'years old!')
]
log('hello', name)
log('we do not know your age anymore because it got deleted when the wrap finished.')
```

A wrap, unless specificifed otherwise, can have a return value. This value can be used as the value of the wrap.

example:

```
const password = [
    let password = input("set password: ")
    while len(password) < 8 [
        log("password must be longer then 8 characters!")
        password = input("set password: ")
    ]
    return password
]

log("your password is", password)
```

If the wrap does not take a return value, then the wrap passes the return value back to a parent wrap.
