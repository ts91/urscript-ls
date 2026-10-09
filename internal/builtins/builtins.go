package builtins

// Builtin describes a built-in URScript function.
type Builtin struct {
	Name      string
	Signature string
	Doc       string
}

// Registry holds the built-in function definitions.
type Registry struct {
	funcs map[string]Builtin
}

// NewRegistry creates a registry pre-populated with standard UR built-ins.
func NewRegistry() *Registry {
	r := &Registry{funcs: make(map[string]Builtin)}
	for _, b := range defaultBuiltins {
		r.funcs[b.Name] = b
	}
	return r
}

// Lookup returns the built-in with the given name, if any.
func (r *Registry) Lookup(name string) (Builtin, bool) {
	b, ok := r.funcs[name]
	return b, ok
}

var defaultBuiltins = []Builtin{
	// Motion
	{
		Name:      "movej",
		Signature: "movej(q, a=1.4, v=1.05, t=0, r=0)",
		Doc:       "Move to position (linear in joint-space). `q` is joint positions or pose, `a` is joint acceleration [rad/s²], `v` is joint speed [rad/s], `t` is time [s] (0 = not specified), `r` is blend radius [m].",
	},
	{
		Name:      "movel",
		Signature: "movel(pose, a=1.2, v=0.25, t=0, r=0)",
		Doc:       "Move to position (linear in tool-space). `pose` is target pose, `a` is tool acceleration [m/s²], `v` is tool speed [m/s], `t` is time [s], `r` is blend radius [m].",
	},
	{
		Name:      "movep",
		Signature: "movep(pose, a=1.2, v=0.25, r=0)",
		Doc:       "Move process: blend circular in tool-space. `pose` is target, `a` is tool acceleration [m/s²], `v` is tool speed [m/s], `r` is blend radius [m].",
	},
	{
		Name:      "movec",
		Signature: "movec(pose_via, pose_to, a=1.2, v=0.25, r=0, mode=0)",
		Doc:       "Move circular via `pose_via` to `pose_to`. `a` is tool acceleration [m/s²], `v` is tool speed [m/s], `r` is blend radius [m], `mode` selects orientation mode.",
	},
	{
		Name:      "servoj",
		Signature: "servoj(q, a, v, t=0.008, lookahead_time=0.1, gain=300)",
		Doc:       "Servo to position in joint-space. Used for online real-time control.",
	},
	{
		Name:      "speedl",
		Signature: "speedl(xd, a=1.2, t=0.0, aRot=a)",
		Doc:       "Accelerate linearly in tool-space. `xd` is tool speed [m/s], `a` is tool acceleration [m/s²], `t` is time [s].",
	},
	{
		Name:      "speedj",
		Signature: "speedj(qd, a, t=0.0)",
		Doc:       "Accelerate in joint-space. `qd` is joint speeds [rad/s], `a` is joint acceleration [rad/s²], `t` is time [s].",
	},
	{
		Name:      "stopj",
		Signature: "stopj(a)",
		Doc:       "Decelerate joint speeds to zero. `a` is joint acceleration [rad/s²].",
	},
	{
		Name:      "stopl",
		Signature: "stopl(a, aRot=a)",
		Doc:       "Decelerate tool speed to zero. `a` is tool deceleration [m/s²].",
	},

	// I/O
	{
		Name:      "set_digital_out",
		Signature: "set_digital_out(n, b)",
		Doc:       "Set digital output signal `n` to value `b` (True/False).",
	},
	{
		Name:      "get_digital_in",
		Signature: "get_digital_in(n)",
		Doc:       "Read digital input signal `n`. Returns True or False.",
	},
	{
		Name:      "set_analog_out",
		Signature: "set_analog_out(n, f)",
		Doc:       "Set analog output `n` to value `f` [mA or V].",
	},
	{
		Name:      "get_analog_in",
		Signature: "get_analog_in(n)",
		Doc:       "Read analog input `n`. Returns a float.",
	},
	{
		Name:      "set_tool_digital_out",
		Signature: "set_tool_digital_out(n, b)",
		Doc:       "Set tool digital output `n` to value `b` (True/False).",
	},
	{
		Name:      "get_tool_digital_in",
		Signature: "get_tool_digital_in(n)",
		Doc:       "Read tool digital input `n`. Returns True or False.",
	},
	{
		Name:      "set_standard_analog_out",
		Signature: "set_standard_analog_out(n, f)",
		Doc:       "Set standard analog output `n` to value `f`.",
	},
	{
		Name:      "get_standard_analog_in",
		Signature: "get_standard_analog_in(n)",
		Doc:       "Read standard analog input `n`. Returns a float.",
	},
	{
		Name:      "set_tool_voltage",
		Signature: "set_tool_voltage(voltage)",
		Doc:       "Set tool voltage. `voltage` is 0, 12, or 24 [V].",
	},

	// Robot status
	{
		Name:      "get_actual_joint_positions",
		Signature: "get_actual_joint_positions()",
		Doc:       "Returns the current joint angular positions as a 6-element array [rad].",
	},
	{
		Name:      "get_actual_tcp_pose",
		Signature: "get_actual_tcp_pose()",
		Doc:       "Returns the current TCP pose as p[x, y, z, rx, ry, rz].",
	},
	{
		Name:      "get_actual_joint_speeds",
		Signature: "get_actual_joint_speeds()",
		Doc:       "Returns the current joint angular velocities [rad/s].",
	},
	{
		Name:      "get_actual_tcp_speed",
		Signature: "get_actual_tcp_speed()",
		Doc:       "Returns the current TCP speed as a 6-element vector [m/s, rad/s].",
	},
	{
		Name:      "get_joint_torques",
		Signature: "get_joint_torques()",
		Doc:       "Returns the current joint torques [Nm].",
	},
	{
		Name:      "get_target_joint_positions",
		Signature: "get_target_joint_positions()",
		Doc:       "Returns the target joint positions [rad].",
	},
	{
		Name:      "get_target_tcp_pose",
		Signature: "get_target_tcp_pose()",
		Doc:       "Returns the target TCP pose as p[x, y, z, rx, ry, rz].",
	},
	{
		Name:      "is_steady",
		Signature: "is_steady()",
		Doc:       "Returns True when the robot is fully at rest.",
	},
	{
		Name:      "get_controller_temp",
		Signature: "get_controller_temp()",
		Doc:       "Returns the controller temperature [°C].",
	},
	{
		Name:      "get_joint_temp",
		Signature: "get_joint_temp(j)",
		Doc:       "Returns the temperature of joint `j` [°C].",
	},
	{
		Name:      "powerdown",
		Signature: "powerdown()",
		Doc:       "Shut down the robot, and when possible, the controller.",
	},

	// Utility
	{
		Name:      "sleep",
		Signature: "sleep(t)",
		Doc:       "Pause program execution for `t` seconds.",
	},
	{
		Name:      "textmsg",
		Signature: "textmsg(s1, s2=\"\")",
		Doc:       "Send text message to the log. `s1` is the main message, `s2` is a title.",
	},
	{
		Name:      "popup",
		Signature: "popup(s, title=\"Popup\", warning=False, error=False, blocking=False)",
		Doc:       "Display a popup message on the teach pendant.",
	},
	{
		Name:      "sync",
		Signature: "sync()",
		Doc:       "Uses up the remaining physical time in the current frame.",
	},

	// Math
	{
		Name:      "norm",
		Signature: "norm(a)",
		Doc:       "Returns the Euclidean norm of a pose or vector.",
	},
	{
		Name:      "point_dist",
		Signature: "point_dist(p_from, p_to)",
		Doc:       "Returns the distance between two points (ignoring rotation).",
	},
	{
		Name:      "pose_add",
		Signature: "pose_add(p_1, p_2)",
		Doc:       "Add two poses together.",
	},
	{
		Name:      "pose_dist",
		Signature: "pose_dist(p_from, p_to)",
		Doc:       "Returns the distance between two poses.",
	},
	{
		Name:      "pose_inv",
		Signature: "pose_inv(p)",
		Doc:       "Returns the inverse of a pose transformation.",
	},
	{
		Name:      "pose_trans",
		Signature: "pose_trans(p_from, p_from_to)",
		Doc:       "Returns the pose resulting from `p_from` transformed by `p_from_to`.",
	},
	{
		Name:      "rotvec2rpy",
		Signature: "rotvec2rpy(rotation_vector)",
		Doc:       "Convert a rotation vector [rx, ry, rz] to roll-pitch-yaw [r, p, y].",
	},
	{
		Name:      "rpy2rotvec",
		Signature: "rpy2rotvec(rpy_vector)",
		Doc:       "Convert roll-pitch-yaw [r, p, y] to a rotation vector [rx, ry, rz].",
	},
	{
		Name:      "interpolate_pose",
		Signature: "interpolate_pose(p_from, p_to, alpha)",
		Doc:       "Linear interpolation between two poses. `alpha` ranges from 0 to 1.",
	},
	{
		Name:      "d2r",
		Signature: "d2r(d)",
		Doc:       "Convert degrees to radians.",
	},
	{
		Name:      "r2d",
		Signature: "r2d(r)",
		Doc:       "Convert radians to degrees.",
	},
	{
		Name:      "get_list_length",
		Signature: "get_list_length(v)",
		Doc:       "Returns the length of a list or array.",
	},
	{
		Name:      "length",
		Signature: "length(v)",
		Doc:       "Returns the length of a list, string, or array.",
	},
	{
		Name:      "append",
		Signature: "append(list, value)",
		Doc:       "Append `value` to the end of `list`. Returns the extended list.",
	},
	{
		Name:      "str_at",
		Signature: "str_at(src, index)",
		Doc:       "Returns the character at `index` in the string `src`.",
	},
	{
		Name:      "str_cat",
		Signature: "str_cat(op1, op2)",
		Doc:       "Concatenate two values as strings.",
	},
	{
		Name:      "str_find",
		Signature: "str_find(src, target)",
		Doc:       "Find first occurrence of `target` in `src`. Returns index or -1.",
	},
	{
		Name:      "str_len",
		Signature: "str_len(s)",
		Doc:       "Returns the length of string `s`.",
	},
	{
		Name:      "str_sub",
		Signature: "str_sub(src, index, length)",
		Doc:       "Returns a substring of `src` starting at `index` with given `length`.",
	},
	{
		Name:      "to_num",
		Signature: "to_num(str)",
		Doc:       "Convert string to number.",
	},
	{
		Name:      "to_str",
		Signature: "to_str(val)",
		Doc:       "Convert value to string representation.",
	},
	{
		Name:      "floor",
		Signature: "floor(val)",
		Doc:       "Round `val` down to nearest integer.",
	},
	{
		Name:      "ceil",
		Signature: "ceil(val)",
		Doc:       "Round `val` up to nearest integer.",
	},
	{
		Name:      "abs",
		Signature: "abs(val)",
		Doc:       "Returns the absolute value of `val`.",
	},
	{
		Name:      "acos",
		Signature: "acos(val)",
		Doc:       "Returns the arc cosine of `val` [rad].",
	},
	{
		Name:      "asin",
		Signature: "asin(val)",
		Doc:       "Returns the arc sine of `val` [rad].",
	},
	{
		Name:      "atan",
		Signature: "atan(val)",
		Doc:       "Returns the arc tangent of `val` [rad].",
	},
	{
		Name:      "atan2",
		Signature: "atan2(y, x)",
		Doc:       "Returns the arc tangent of `y/x` [rad], using signs to determine quadrant.",
	},
	{
		Name:      "cos",
		Signature: "cos(val)",
		Doc:       "Returns the cosine of `val` [rad].",
	},
	{
		Name:      "sin",
		Signature: "sin(val)",
		Doc:       "Returns the sine of `val` [rad].",
	},
	{
		Name:      "tan",
		Signature: "tan(val)",
		Doc:       "Returns the tangent of `val` [rad].",
	},
	{
		Name:      "sqrt",
		Signature: "sqrt(val)",
		Doc:       "Returns the square root of `val`.",
	},
	{
		Name:      "pow",
		Signature: "pow(base, exp)",
		Doc:       "Returns `base` raised to the power `exp`.",
	},
	{
		Name:      "log",
		Signature: "log(val)",
		Doc:       "Returns the natural logarithm of `val`.",
	},
	{
		Name:      "random",
		Signature: "random()",
		Doc:       "Returns a pseudo-random number between 0 and 1.",
	},

	// Communication
	{
		Name:      "socket_open",
		Signature: "socket_open(address, port, socket_name=\"socket_0\")",
		Doc:       "Open a TCP/IP socket connection to `address`:`port`.",
	},
	{
		Name:      "socket_close",
		Signature: "socket_close(socket_name=\"socket_0\")",
		Doc:       "Close the socket with the given name.",
	},
	{
		Name:      "socket_send_string",
		Signature: "socket_send_string(str, socket_name=\"socket_0\")",
		Doc:       "Send `str` through the socket.",
	},
	{
		Name:      "socket_send_line",
		Signature: "socket_send_line(str, socket_name=\"socket_0\")",
		Doc:       "Send `str` followed by a newline through the socket.",
	},
	{
		Name:      "socket_send_byte",
		Signature: "socket_send_byte(value, socket_name=\"socket_0\")",
		Doc:       "Send a single byte through the socket.",
	},
	{
		Name:      "socket_read_string",
		Signature: "socket_read_string(socket_name=\"socket_0\", prefix=\"\", suffix=\"\", interpret_escape=False)",
		Doc:       "Read a string from the socket.",
	},
	{
		Name:      "socket_read_byte_list",
		Signature: "socket_read_byte_list(number, socket_name=\"socket_0\", timeout=2)",
		Doc:       "Read `number` bytes from the socket. Returns a list.",
	},
	{
		Name:      "socket_read_ascii_float",
		Signature: "socket_read_ascii_float(number, socket_name=\"socket_0\", timeout=2)",
		Doc:       "Read `number` ASCII-encoded floats from the socket.",
	},
	{
		Name:      "socket_read_binary_integer",
		Signature: "socket_read_binary_integer(number, socket_name=\"socket_0\", timeout=2)",
		Doc:       "Read `number` binary-encoded integers from the socket.",
	},

	// Force/Torque
	{
		Name:      "force",
		Signature: "force()",
		Doc:       "Returns the norm of the current TCP force [N].",
	},
	{
		Name:      "get_tcp_force",
		Signature: "get_tcp_force()",
		Doc:       "Returns the current TCP wrench (force and torque).",
	},
	{
		Name:      "force_mode",
		Signature: "force_mode(task_frame, selection_vector, wrench, type, limits)",
		Doc:       "Set robot in force mode. The robot complies in the specified directions.",
	},
	{
		Name:      "end_force_mode",
		Signature: "end_force_mode()",
		Doc:       "Exit force mode.",
	},
	{
		Name:      "zero_ftsensor",
		Signature: "zero_ftsensor()",
		Doc:       "Zero the force/torque sensor.",
	},

	// TCP
	{
		Name:      "set_tcp",
		Signature: "set_tcp(pose)",
		Doc:       "Set the tool center point (TCP) offset.",
	},
	{
		Name:      "set_payload",
		Signature: "set_payload(m, cog=[0,0,0])",
		Doc:       "Set payload mass `m` [kg] and center of gravity `cog` [m].",
	},
	{
		Name:      "set_gravity",
		Signature: "set_gravity(d)",
		Doc:       "Set the direction of gravity as a 3D vector (typically [0,0,9.82]).",
	},

	// Conveyor tracking
	{
		Name:      "conveyor_pulse_decode",
		Signature: "conveyor_pulse_decode(type, A, B)",
		Doc:       "Set up conveyor encoder. `type`: 0=no encoder, 1=quadrature, 2=rising, 3=falling.",
	},

	// Threads
	{
		Name:      "run",
		Signature: "run(thread_fn)",
		Doc:       "Start a thread executing `thread_fn`.",
	},
	{
		Name:      "kill",
		Signature: "kill(thread_handle)",
		Doc:       "Kill a running thread.",
	},
	{
		Name:      "enter_critical",
		Signature: "enter_critical()",
		Doc:       "Enter a critical section (blocks thread switching).",
	},
	{
		Name:      "exit_critical",
		Signature: "exit_critical()",
		Doc:       "Exit a critical section (re-enables thread switching).",
	},

	// Pose constructor
	{
		Name:      "p",
		Signature: "p[x, y, z, rx, ry, rz]",
		Doc:       "Pose literal constructor. Represents a position and orientation in space.",
	},

	// Types
	{
		Name:      "is_number",
		Signature: "is_number(val)",
		Doc:       "Returns True if `val` is a number.",
	},
	{
		Name:      "is_string",
		Signature: "is_string(val)",
		Doc:       "Returns True if `val` is a string.",
	},
	{
		Name:      "is_list",
		Signature: "is_list(val)",
		Doc:       "Returns True if `val` is a list.",
	},
	{
		Name:      "is_pose",
		Signature: "is_pose(val)",
		Doc:       "Returns True if `val` is a pose.",
	},

	// Misc
	{
		Name:      "set_flag",
		Signature: "set_flag(n, b)",
		Doc:       "Set flag `n` to boolean value `b`.",
	},
	{
		Name:      "get_flag",
		Signature: "get_flag(n)",
		Doc:       "Returns the value of flag `n`.",
	},
	{
		Name:      "set_pos",
		Signature: "set_pos(name, pos)",
		Doc:       "Store a position `pos` with name `name`.",
	},
	{
		Name:      "get_pos",
		Signature: "get_pos(name)",
		Doc:       "Retrieve the stored position with name `name`.",
	},
	{
		Name:      "write_output_integer_register",
		Signature: "write_output_integer_register(address, value)",
		Doc:       "Write `value` to the integer output register at `address`.",
	},
	{
		Name:      "read_input_integer_register",
		Signature: "read_input_integer_register(address)",
		Doc:       "Read the integer input register at `address`.",
	},
	{
		Name:      "write_output_float_register",
		Signature: "write_output_float_register(address, value)",
		Doc:       "Write `value` to the float output register at `address`.",
	},
	{
		Name:      "read_input_float_register",
		Signature: "read_input_float_register(address)",
		Doc:       "Read the float input register at `address`.",
	},
	{
		Name:      "write_output_boolean_register",
		Signature: "write_output_boolean_register(address, value)",
		Doc:       "Write `value` to the boolean output register at `address`.",
	},
	{
		Name:      "read_input_boolean_register",
		Signature: "read_input_boolean_register(address)",
		Doc:       "Read the boolean input register at `address`.",
	},
	{
		Name:      "modbus_add_signal",
		Signature: "modbus_add_signal(IP, slave_number, signal_address, signal_type, signal_name)",
		Doc:       "Add a Modbus signal.",
	},
	{
		Name:      "modbus_get_signal_status",
		Signature: "modbus_get_signal_status(signal_name, is_secondary_program)",
		Doc:       "Returns the value of the named Modbus signal.",
	},
	{
		Name:      "modbus_set_output_signal",
		Signature: "modbus_set_output_signal(signal_name, value, is_secondary_program)",
		Doc:       "Set the value of a Modbus output signal.",
	},
}
