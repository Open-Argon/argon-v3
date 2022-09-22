print('1. solve hypotenuse')
print('2. solve short side')
print('3. quit')

while True then
  choice < - input('choice: ')
   if choice == 1 then
     a = input('a: ')
      b = input('b: ')
       c = √(a*a+b*b)
        print('c:', c)
    elif choice == 2 then
     b = input('b: ')
      c = input('c: ')
       a = √(c*c-b*b)
        print('a:', a)
    else then
     quit()
    endif
endwhile
