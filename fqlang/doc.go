// Package fql implements a tokenizer for FQL.
//
// Tokenization is done by creating a Tokenizer for a string that contains the fql.
//
//	z := fql.NewTokenizer(fqlstring)
//
// Given a Tokenizer z, the FQL is tokenized by repeatedly calling z.Next(),
// which parses the next token and returns it, or an error. If the returned
// token is fql.RULE, then calling z.Rule() will return the parsed *fql.Rule
// value that is associated with that token.
//
//	for {
//		tok, err := z.Next()
//		if err != nil {
//			if err != fql.EOF {
//				return err
//			}
//			break
//		}
//
//		switch tok {
//		case fql.LPAREN:
//			// ...
//		case fql.RPAREN:
//			// ...
//		case fql.AND:
//			// ...
//		case fql.OR:
//			// ...
//		case fql.RULE:
//			r := z.Rule()
//			// ...
//		}
//	}
package fqlang
