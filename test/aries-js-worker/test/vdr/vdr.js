/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {newAries, newAriesREST} from "../common.js"
import {environment} from "../environment.js";

const agentControllerApiUrl = `${environment.HTTP_SCHEME}://${environment.USER_HOST}:${environment.USER_API_PORT}`

// did document
const didName = "alice-did"
const didID = `${environment.DID_ID}`


// scenarios
describe("VDR", function () {
    let agents
    before(async () => {
        let a = await newAries('vdr-demo','vdr-demo-agent', ["sidetree@http://localhost:48326/sidetree/0.0.1/identifiers"])
        let b = await newAriesREST(agentControllerApiUrl)

        agents = [a, b]
    })

    after(async () => {
        await agents[0].destroy()
        await agents[1].destroy()
    })

    it("Alice stores the did generated by her", async function () {
        for (let i = 0; i < agents.length; i++) {
            let resp;
            for (let x = 0; x < 10; x++) {
                try {
                    let id=didID
                    if (i===1){
                        id=window.btoa(didID)
                    }

                    resp = await agents[i].vdr.resolveDID({id: id})

                    break
                } catch (e) {
                    if (!e.message.includes("DID does not exist")) {
                        assert.fail(e.message);
                    }
                    await new Promise(r => setTimeout(r, 1000));
                    console.warn(e.message)
                }
            }
            await agents[i].vdr.saveDID({
                name: didName,
                did: resp.did
            })
        }
    })

    it("Alice retrieves the did from store",  async function () {
        let errors = "Alice didn't retrieve did: ";

        for (let i = 0; i < agents.length; i++) {
            try {
                await agents[i].vdr.getDID({id: didID})
            }catch (e) {
                errors += e.message;
                continue
            }

            return
        }

        throw new Error(errors)
    })

    it("Alice validates that she has the did", async function () {
        for (let i = 0; i < agents.length; i++) {
            let resp;
            try {
                resp = await agents[i].vdr.getDIDRecords()
            }catch (e) {
                continue
            }

            for (let j = 0 ;j<resp.result.length;j++){
                if (didID === resp.result[j].id){
                    assert.equal(didName, resp.result[j].name)
                    return
                }
            }
        }

        throw new Error("Alice doesn't have did")
    })
})
