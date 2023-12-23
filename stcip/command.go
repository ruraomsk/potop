package stcip

// Команды:
// вызов фазы - 1.3.6.1.4.1.1618.3.7.2.11.1.0
// SwarcoUTCTrafftechPhaseCommand OBJECT-TYPE
// SYNTAX Unsigned32
// ACCESS read-write
// STATUS current
// DESCRIPTION
// "Commands the controller to go to the
// specified phase. A phase of 0 means no phase commanded."
// ::= { SwarcoUTCTrafftechPhaseObjs 1 }

// вызов плана - 1.3.6.1.4.1.1618.3.7.2.2.1.0
// SwarcoUTCTrafftechPlanCommand OBJECT-TYPE
// SYNTAX Unsigned32
// ACCESS read-write
// STATUS current
// DESCRIPTION
// "Command the controller to enter the given
// trafiic plan number.
// traffic
// 0 means automatic plan selection otherwise a
// plan number to be forced."
// ::= { SwarcoUTCTrafftechPlanCommandObjs 1 }

// вкл/выкл ЖМ - 1.3.6.1.4.1.1618.3.2.2.1.1.0
// SwarcoUTCCommandFlash OBJECT-TYPE
// SYNTAX CommandFlash
// ACCESS read-write
// STATUS current
// DESCRIPTION
// "The flash command."
// ::= { SwarcoUTCCommandFlashObjs 1 }

// -- 1.3.6.1.4.1.1618.3.2.2.2.1
// SwarcoUTCCommandDark OBJECT-TYPE
// SYNTAX CommandDark
// ACCESS read-write
// STATUS current
// DESCRIPTION
// "The dark command."
// ::= { SwarcoUTCCommandDarkObjs 1 }
