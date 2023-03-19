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

`(let or nothing) {NAME} = {ANY}`

let variables will set variables. at the end of a opperation (e.g. if, while, or function) drop the
variable stack.

having no verb at the start suggests the program would like to edit a variables from a different stack.
if there is no variable found in the other stacks, then it sets one in the current stack as a let.

setting variables returns the value, which can be used.

example:

```javascript
if ((x = 10) > 5) do
    term.log(x, 'is bigger than 5')
```

---

## functions

`(let or nothing) {NAME}({PARAMS}) = {CODE}`

the starting verb follows the rules set by [set variable](#set-variable).

function can be called by using its name followed by brackets with the params.

example:

```javascript
f(x) = x^2 + 2*x + 1
term.log('f(10) =', f(10))
```

output:

```javascript
f(10) = 121
```

if the function does not return, then the value returned is `null`

---

## wrap

```
do
   {CODE}
```

a wrap encloses code indented in after the word `do` its used to create a new scope, so variables set from
inside the wraps scope are deleted once the wrap is finished.

example:

```javascript
let name = null

do
    name = input('name: ')
    let age = input('age: ')
    log('you are', age, 'years old!')

term.log('hello', name)
term.log('we do not know your age anymore because it got deleted when the wrap finished.')
```

A wrap, unless specificifed otherwise, can have a return value. This value can be used as the value of the wrap.

example:

```javascript
let password = do
    let password = passwordInput("set password: ")
    while (password.length < 8) do
        term.log("password must be longer then 8 characters!")
        password = passwordInput("set password: ")
    return password

term.log("your password is", password)
```

If the wrap does not take a return value, then the wrap passes the return value back to a parent wrap.

## Comments
`//{COMMENT}`

Comments allow the programmer to write a message into their code, without the message being processed by the computer.