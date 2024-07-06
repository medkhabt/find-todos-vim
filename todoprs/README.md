## Grammar 
### parsing rg output 
$`R \ => \ P_0 \ l \ E \ T`$ 

$`P_0 \ => \ c_0 \ P_1`$

$`P_1 \ => \ c_0 \ P_1 \ | \ / \ c_0 \ P_1 \ | \ . \ P_1 \ | \ \varepsilon`$

$`E \ => \ td \ T \ | \ c_0 \ E \ | \ cmt \ E`$

$`T \ => \ c_0 \ T \ | \ sp \ T \ | \ \varepsilon`$

$c_0$ terminal is [a-z]|[A-Z]|[0-9]|-|_$

$l$ terminal is the line_number token

$cmt$ terminal is the leading single line comment token.

$sp$ terminal is the .|/|#/$/ /

I might not consider the special charachters for now. 

