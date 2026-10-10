package builtins

func urmathBuiltins() []Builtin {
	return []Builtin{
		{
			Name:      "acos",
			Signature: "acos(f)",
			Module:    "urmath",
			Doc: "Returns the arc cosine of f in radians. Raises error if f is outside [-1, 1].\n" +
				"@param f float Value\n" +
				"@returns float Arc cosine in radians",
		},
		{
			Name:      "asin",
			Signature: "asin(f)",
			Module:    "urmath",
			Doc: "Returns the arc sine of f in radians. Raises error if f is outside [-1, 1].\n" +
				"@param f float Value\n" +
				"@returns float Arc sine in radians",
		},
		{
			Name:      "atan",
			Signature: "atan(f)",
			Module:    "urmath",
			Doc: "Returns the arc tangent of f in radians.\n" +
				"@param f float Value\n" +
				"@returns float Arc tangent in radians",
		},
		{
			Name:      "atan2",
			Signature: "atan2(x, y)",
			Module:    "urmath",
			Doc: "Returns the arc tangent of x/y in radians, using signs to determine the quadrant.\n" +
				"@param x float Numerator\n" +
				"@param y float Denominator\n" +
				"@returns float Arc tangent in radians",
		},
		{
			Name:      "binary_list_to_integer",
			Signature: "binary_list_to_integer(l)",
			Module:    "urmath",
			Doc: "Convert a list of booleans to a signed integer (index 0 is LSB).\n" +
				"@param l list List of bools representing binary digits\n" +
				"@returns int Signed integer value (uses first 32 bits if list is longer)",
		},
		{
			Name:      "ceil",
			Signature: "ceil(f)",
			Module:    "urmath",
			Doc: "Returns the smallest integer not less than f.\n" +
				"@param f float Value to round up\n" +
				"@returns int Rounded integer",
		},
		{
			Name:      "cos",
			Signature: "cos(f)",
			Module:    "urmath",
			Doc: "Returns the cosine of f.\n" +
				"@param f float Angle in radians\n" +
				"@returns float Cosine value",
		},
		{
			Name:      "d2r",
			Signature: "d2r(d)",
			Module:    "urmath",
			Doc: "Convert degrees to radians.\n" +
				"@param d float Angle in degrees\n" +
				"@returns float Angle in radians",
		},
		{
			Name:      "floor",
			Signature: "floor(f)",
			Module:    "urmath",
			Doc: "Returns the largest integer not greater than f.\n" +
				"@param f float Value to round down\n" +
				"@returns int Rounded integer",
		},
		{
			Name:      "make_list",
			Signature: "make_list(length, initial_value, capacity=length)",
			Module:    "urmath",
			Doc: "Create a new list with the given length, initial value, and optional capacity.\n" +
				"@param length int Number of elements to initialize\n" +
				"@param initial_value any Initial value (also sets the list element type)\n" +
				"@param capacity int Maximum number of elements (defaults to length)\n" +
				"@returns list The new list",
		},
		{
			Name:      "get_list_length",
			Signature: "get_list_length(v)",
			Module:    "urmath",
			Doc: "Returns the number of entries in a list.\n" +
				"@param v list A list variable\n" +
				"@returns int Length of the list",
		},
		{
			Name:      "integer_to_binary_list",
			Signature: "integer_to_binary_list(x)",
			Module:    "urmath",
			Doc: "Convert a signed integer to a list of 32 booleans (index 0 is LSB).\n" +
				"@param x int Integer value to convert\n" +
				"@returns list List of 32 bools",
		},
		{
			Name:      "interpolate_pose",
			Signature: "interpolate_pose(p_from, p_to, alpha)",
			Module:    "urmath",
			Doc: "Linear interpolation of tool position and geodetic orientation.\n" +
				"alpha=0 returns p_from, alpha=1 returns p_to. Values outside [0,1] extrapolate.\n" +
				"@param p_from pose Start pose\n" +
				"@param p_to pose End pose\n" +
				"@param alpha float Interpolation factor\n" +
				"@returns pose Interpolated pose",
		},
		{
			Name:      "inv",
			Signature: "inv(m)",
			Module:    "urmath",
			Doc: "Returns the inverse of a square matrix or pose.\n" +
				"@param m any Matrix or pose (must be square and non-singular)\n" +
				"@returns any Inverse matrix or pose",
		},
		{
			Name:      "length",
			Signature: "length(v)",
			Module:    "urmath",
			Doc: "Returns the length of a list or string.\n" +
				"@param v any List or string\n" +
				"@returns int Number of entries or characters",
		},
		{
			Name:      "log",
			Signature: "log(b, f)",
			Module:    "urmath",
			Doc: "Returns the logarithm of f to base b. Raises error if b or f is negative, or b is 1.\n" +
				"@param b float Base\n" +
				"@param f float Value\n" +
				"@returns float Logarithm result",
		},
		{
			Name:      "norm",
			Signature: "norm(a)",
			Module:    "urmath",
			Doc: "Returns the norm of the argument.\n" +
				"For poses and lists: Euclidean norm. For floats: absolute value. For ints: absolute value.\n" +
				"@param a any Pose, float, int, or list of numbers\n" +
				"@returns float Norm value",
		},
		{
			Name:      "normalize",
			Signature: "normalize(v)",
			Module:    "urmath",
			Doc: "Returns the unit vector in the direction of v. Raises error if all elements are zero.\n" +
				"@param v list List of floats\n" +
				"@returns list Normalized vector",
		},
		{
			Name:      "point_dist",
			Signature: "point_dist(p_from, p_to)",
			Module:    "urmath",
			Doc: "Returns the distance between two points, ignoring orientation.\n" +
				"@param p_from pose First point\n" +
				"@param p_to pose Second point\n" +
				"@returns float Distance in meters",
		},
		{
			Name:      "pose_add",
			Signature: "pose_add(p_1, p_2)",
			Module:    "urmath",
			Doc: "Add two poses: positions are summed, rotations are multiplied.\n" +
				"@param p_1 pose First pose\n" +
				"@param p_2 pose Second pose\n" +
				"@returns pose Sum pose",
		},
		{
			Name:      "pose_dist",
			Signature: "pose_dist(p_from, p_to)",
			Module:    "urmath",
			Doc: "Returns the distance between two poses, including rotational difference.\n" +
				"@param p_from pose First pose\n" +
				"@param p_to pose Second pose\n" +
				"@returns float Distance",
		},
		{
			Name:      "pose_inv",
			Signature: "pose_inv(p_from)",
			Module:    "urmath",
			Doc: "Returns the inverse of a pose transformation.\n" +
				"@param p_from pose Pose to invert\n" +
				"@returns pose Inverse pose",
		},
		{
			Name:      "pose_sub",
			Signature: "pose_sub(p_to, p_from)",
			Module:    "urmath",
			Doc: "Subtract two poses.\n" +
				"@param p_to pose Target pose\n" +
				"@param p_from pose Base pose\n" +
				"@returns pose Difference pose",
		},
		{
			Name:      "pose_trans",
			Signature: "pose_trans(p_from, p_from_to)",
			Module:    "urmath",
			Doc: "Pose transformation: apply p_from_to relative to p_from.\n" +
				"Equivalent to matrix multiplication T_world_to = T_world_from * T_from_to.\n" +
				"@param p_from pose Starting pose\n" +
				"@param p_from_to pose Relative transformation\n" +
				"@returns pose Resulting pose",
		},
		{
			Name:      "pow",
			Signature: "pow(base, exponent)",
			Module:    "urmath",
			Doc: "Returns base raised to the power of exponent.\n" +
				"@param base float Base value\n" +
				"@param exponent float Exponent value\n" +
				"@returns float Result",
		},
		{
			Name:      "r2d",
			Signature: "r2d(r)",
			Module:    "urmath",
			Doc: "Convert radians to degrees.\n" +
				"@param r float Angle in radians\n" +
				"@returns float Angle in degrees",
		},
		{
			Name:      "random",
			Signature: "random()",
			Module:    "urmath",
			Doc: "Returns a pseudo-random number between 0 and 1.\n" +
				"@returns float Random value in [0, 1]",
		},
		{
			Name:      "rotvec2rpy",
			Signature: "rotvec2rpy(rotation_vector)",
			Module:    "urmath",
			Doc: "Convert a rotation vector (axis-angle) to roll-pitch-yaw angles.\n" +
				"@param rotation_vector list Rotation vector [rx, ry, rz] in radians\n" +
				"@returns list RPY vector [roll, pitch, yaw] in radians",
		},
		{
			Name:      "rpy2rotvec",
			Signature: "rpy2rotvec(rpy_vector)",
			Module:    "urmath",
			Doc: "Convert roll-pitch-yaw angles to a rotation vector (axis-angle).\n" +
				"@param rpy_vector list RPY vector [roll, pitch, yaw] in radians\n" +
				"@returns list Rotation vector [rx, ry, rz] in radians",
		},
		{
			Name:      "sin",
			Signature: "sin(f)",
			Module:    "urmath",
			Doc: "Returns the sine of f.\n" +
				"@param f float Angle in radians\n" +
				"@returns float Sine value",
		},
		{
			Name:      "size",
			Signature: "size(v)",
			Module:    "urmath",
			Doc: "Returns the size of a matrix, list, or string.\n" +
				"For lists/strings: integer length. For matrices: [rows, cols].\n" +
				"@param v any Matrix, list, or string\n" +
				"@returns any Integer length or [rows, cols] list",
		},
		{
			Name:      "sqrt",
			Signature: "sqrt(f)",
			Module:    "urmath",
			Doc: "Returns the square root of f. Raises error if f is negative.\n" +
				"@param f float Value\n" +
				"@returns float Square root",
		},
		{
			Name:      "tan",
			Signature: "tan(f)",
			Module:    "urmath",
			Doc: "Returns the tangent of f.\n" +
				"@param f float Angle in radians\n" +
				"@returns float Tangent value",
		},
		{
			Name:      "transpose",
			Signature: "transpose(m)",
			Module:    "urmath",
			Doc: "Returns the transpose of a matrix or array.\n" +
				"@param m any Matrix or array\n" +
				"@returns any Transposed matrix or array",
		},
		{
			Name:      "wrench_trans",
			Signature: "wrench_trans(T_from_to, w_from)",
			Module:    "urmath",
			Doc: "Transform a wrench to a different point of view.\n" +
				"@param T_from_to pose Transformation to the new reference frame\n" +
				"@param w_from list Wrench [Fx, Fy, Fz, Mx, My, Mz]\n" +
				"@returns list Transformed wrench [Fx, Fy, Fz, Mx, My, Mz]",
		},
		{
			Name:      "abs",
			Signature: "abs(val)",
			Module:    "urmath",
			Doc: "Returns the absolute value.\n" +
				"@param val number Numeric value\n" +
				"@returns number Absolute value",
		},
		{
			Name:      "append",
			Signature: "append(list, value)",
			Module:    "urmath",
			Doc: "Append a value to the end of a list.\n" +
				"@param list list Target list\n" +
				"@param value any Value to append\n" +
				"@returns list Extended list",
		},
		// Pose constructor
		{
			Name:      "p",
			Signature: "p[x, y, z, rx, ry, rz]",
			Module:    "urmath",
			Doc:       "Pose literal constructor. Represents a position and orientation in space.",
		},
		// Type checks
		{
			Name:      "is_number",
			Signature: "is_number(val)",
			Module:    "urmath",
			Doc: "Check if a value is a number.\n" +
				"@param val any Value to check\n" +
				"@returns bool True if val is a number",
		},
		{
			Name:      "is_string",
			Signature: "is_string(val)",
			Module:    "urmath",
			Doc: "Check if a value is a string.\n" +
				"@param val any Value to check\n" +
				"@returns bool True if val is a string",
		},
		{
			Name:      "is_list",
			Signature: "is_list(val)",
			Module:    "urmath",
			Doc: "Check if a value is a list.\n" +
				"@param val any Value to check\n" +
				"@returns bool True if val is a list",
		},
		{
			Name:      "is_pose",
			Signature: "is_pose(val)",
			Module:    "urmath",
			Doc: "Check if a value is a pose.\n" +
				"@param val any Value to check\n" +
				"@returns bool True if val is a pose",
		},
	}
}
