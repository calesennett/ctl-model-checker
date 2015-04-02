# ctl-model-checker
A basic model checker for CTL.

## FSM File Specification
The .fsm file should adhere to these specifications. A blank line should separate each categorized declaration. Two example .fsm files are included.

### Number of States
    STATES n

### Initial States
Each initial state should be on a separate line:

    INIT
    0
    1

### Transitions
Each arc should be on a separate line:

    ARCS
    0:1
    1:1
    2:0

### Labels
Each label should be specified individually with each desired state below:

    LABEL f
    1
    2

### Properties
Each property should be specified on a separate line:

    (EX (EG f))
    (EG f)

## Usage
    ./scanner < fsm.fsm
