The purpose of this plugin is to provide an alternative for viewing payments whey using multi part payments.

This plugin takes no additional options so to use, start `lightningd` with `--plugin` option

`lightningd --plugin=path/to/plugin`

`mpp_payments` rpc command is available, it takes no parameters.  This command groups payments my payment hash, so each payment will only be listed once.  Non mpp payments will appear as in the `listsendpays` command, mpp payments will display the the total of all payments in the `msatoshi_sent` field.  The `id` field is removed for mpp.  An additional field is added, `parts` which is the total number of parts the payment was split into.