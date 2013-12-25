/*
This Go version 1.2 package is a naive implementation of of the md2 algorithm as defined
in RFC1319 (http://tools.ietf.org/html/rfc1319), for learning purposes only.

This package comes with no warranty and I make no guarantee for its correctness.
Remember that md2 is obsolete, so don't use this for anything else than to
satisfy your curiosity.

I have made no attempt to optimize the code. Instead, the focus has been to
make the code resemble the algorithm as closely as possible. A run o Sum takes
about 25000ns. Relevant tests can be found in the md2_test.go file.

This software is released under The MIT License (MIT) http://opensource.org/licenses/MIT.
See LICENSE file.
*/

package md2

//Permutation of 0..255 constructed from the digits of pi. It gives a
//"random" nonlinear byte substitution operation.

var pi_subst = [256]byte{
	41, 46, 67, 201, 162, 216, 124, 1, 61, 54, 84, 161, 236, 240, 6, 19, 98, 167, 5, 243, 192, 199, 115, 140, 152, 147, 43, 217, 188, 76, 130, 202, 30, 155, 87, 60, 253, 212, 224, 22, 103, 66, 111, 24, 138, 23, 229, 18, 190, 78, 196, 214, 218, 158, 222, 73, 160, 251,
	245, 142, 187, 47, 238, 122, 169, 104, 121, 145, 21, 178, 7, 63, 148, 194, 16, 137, 11, 34, 95, 33, 128, 127, 93, 154, 90, 144, 50, 39, 53, 62, 204, 231, 191, 247, 151, 3, 255, 25, 48, 179, 72, 165, 181, 209, 215, 94, 146, 42, 172, 86, 170, 198, 79, 184, 56, 210,
	150, 164, 125, 182, 118, 252, 107, 226, 156, 116, 4, 241, 69, 157, 112, 89, 100, 113, 135, 32, 134, 91, 207, 101, 230, 45, 168, 2, 27, 96, 37, 173, 174, 176, 185, 246, 28, 70, 97, 105, 52, 64, 126, 15, 85, 71, 163, 35, 221, 81, 175, 58, 195, 92, 249, 206, 186, 197,
	234, 38, 44, 83, 13, 110, 133, 40, 132, 9, 211, 223, 205, 244, 65, 129, 77, 82, 106, 220, 55, 200, 108, 193, 171, 250, 36, 225, 123, 8, 12, 189, 177, 74, 120, 136, 149, 139, 227, 99, 232, 109, 233, 203, 213, 254, 59, 0, 29, 57, 242, 239, 183, 14, 102, 88, 208, 228,
	166, 119, 114, 248, 235, 117, 75, 10, 49, 68, 80, 180, 143, 237, 31, 26, 219, 153, 141, 51, 159, 17, 131, 20,
}

// Section 3.1
// Step 1. Append Padding Bytes
// Adds padding to the end of a byte array 'data' such that the length becomes a multiple of 16

func append_padding_bytes(data []byte) []byte {
	padLen := 16 - len(data)%16
	var padding = make([]byte, padLen)

	/* "i" bytes of value "i" are appended  to the message */
	for i := 0; i < padLen; i++ {
		padding[i] = byte(padLen)
	}

	return append(data, padding[:]...)
}

// Section 3.2
// Step 2. Append Checksum
// This step uses a 256-byte "random" permutation
// constructed from the digits of pi (pi_subst).

func append_checksum(M []byte) []byte {

	// Initialize checksum. Go default value is 0.
	// No initialization necessary.

	var C [16]byte
	for i := 0; i < 16; i++ {
		C[i] = 0
	}

	var L byte = 0
	var c byte

	// Process each 16-word block.
	N := len(M)
	for i := 0; i < N/16; i++ {
		// Checksum block i.
		for j := 0; j < 16; j++ {
			c = M[i*16+j]

			// Notice that an error in the pseudo code in RFC 1319.
			// The errata here: http://www.rfc-editor.org/errata_search.php?rfc=1319
			// corrects;
			// Set C[j] to S[c xor L]
			// to
			// Set C[j] to C[j] xor S[c xor L]
			// Failing to implement this results in incorrect message digests
			// for input strings longer than 16 characters.

			C[j] = C[j] ^ pi_subst[c^L]
			L = C[j]
		}
	}

	return append(M, C[:]...)
}

// Section 3.4
// Step 3 and 4. Process Message in 16-Byte Blocks.

func process_message(M []byte) [16]byte {

	// A 48-byte buffer X is used to compute the message digest.
	// Go default value for bytes is 0. No initialization necessary.

	var X [48]byte

	N := len(M)
	var t int

	// Process each 16-word block.
	for i := 0; i < N/16; i++ {
		// Copy block i into X.
		for j := 0; j < 16; j++ {
			X[16+j] = M[i*16+j]
			X[32+j] = (X[16+j] ^ X[j])
		}

		t = 0

		// Do 18 rounds.
		for j := 0; j < 18; j++ {
			// Round j.
			for k := 0; k < 48; k++ {
				X[k] = (X[k] ^ pi_subst[t])
				t = int(X[k])
			}

			t = (t + j) % 256
		}
	}

	var md [16]byte

	copy(md[:], X[0:16])
	return md
}

// Computes the md2sum for a byte array and returns the 16 byte
// message digest.

func Sum(m []byte) [16]byte {

	// Step 1. Append Padding Bytes
	M := append_padding_bytes(m)

	// Step 2. Append Checksum
	M = append_checksum(M)

	// Step 3 and 4. Initialize MD Buffer and process M in 16-Byte Blocks
	message_digest := process_message(M)

	return message_digest

}
