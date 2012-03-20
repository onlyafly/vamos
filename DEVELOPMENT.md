Vamos Development Notes
=======================

Implementing Tail Recursion
---------------------------

* "CONS Should Not CONS Its Arguments, Part II: Cheney on the M.T.A."
  by Henry Baker
  
  http://home.pipeline.com/~hbaker1/CheneyMTA.html
  
  Describes a technique to implement full tail recursion when
  compiling from Scheme into C.
  
* The Scheme2JS compiler and especially its trampoline implementation has been
  presented:
  http://www-sop.inria.fr/indes/scheme2js/files/tfp2007.pdf
  
* "The 90 Minute Scheme to C compiler" by Marc Feeley
  http://www.iro.umontreal.ca/~boucherd/mslug/meetings/20041020/minutes-en.html
  
  Supports fully optimized proper tail calls, continuations, and (of
  course) full closures, using two important compilation techniques
  for functional languages: closure conversion and CPS-conversion

Implementing call/cc
--------------------

* Description of Scheme2Js's call/cc implementation:
  http://www-sop.inria.fr/indes/scheme2js/files/schemews2007.pdf

Compiler Passes
---------------

* "How Many Passes" by Eric Lippert
  http://blogs.msdn.com/b/ericlippert/archive/2010/02/04/how-many-passes.aspx
