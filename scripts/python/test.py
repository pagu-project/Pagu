import unittest
import public_key
import tempfile
from main import *

class TestPublicKey(unittest.TestCase):
    def test_from_string(self):
        pubStr = "tpublic1pj2aecual5gul3p6epxw0pyulnpxtek3x20r2w5wkguxshfjqx7ecuzmgu0s9g2zkfz3hchasn0u6srwcvprweg20zjy9u5ehj8798lh3jakflaap77u8xqznpkgftg6emkxmg7g3v8052kfjjvtuasqsty299nxn"
        addrAccStr = "tpc1zwmd7fr9muntueu9u0unduy8ycymzgthvu8tgyx"
        addrValStr = "tpc1pwmd7fr9muntueu9u0unduy8ycymzgthvpvm4nm"

        pub = public_key.PublicKey.from_string(pubStr)
        self.assertEqual(pub.account_address().string(), addrAccStr)
        self.assertEqual(pub.validator_address().string(), addrValStr)


if __name__ == '__main__':
    unittest.main()