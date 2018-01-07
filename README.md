# finitefilesystem

Written in an hour because I went out for a walk (always a good idea). I was thinking about how data that has to be stored for users, products, statistics, etc. is made up of two types of values: finite and more-infinite-ish.

The finite values are the booleans, the lists of roles, indexes to other finite datasets; such things. The more-infinite-ish are names, descriptions, pictures, etc. which are obviously not "infinite" in any way (or else you're asking for an exploit) but which at least have an intractable number of possibilities.

I started thinking that it might be useful to split apart these two types of data since they have different storage and usage needs. Most of the time control logic is concerned with checking the finite data while the more "infinite" data is interesting from a processing perspective, it is more the payload. A key thing about the finite data is also that there is a point where it takes less storage to store all permutations (especially if you produce them on the fly) and share them using a hash reference, then to store all the individual objects in each data object you have. That's basically what this is for.

I think there are some interesting advantages to splitting out this kind of data, if only because cracking something open requires us to look at the interior. Some examples:

1. *Versioning* Generally it is not necessary to version user entered data, like names and descriptions, with the same fidelity as we would version data related to control logic during a release for example. Being able to rollback or validate changes when new code, requiring new settings, is released is pretty important. With how compact the stored data can be (since it isn't stored in each object) more old versions can be stored still under the old hashes and you just need to store which hash belonged to which data.
2. *Consistent datasets* Since data is shared it is quicker to update the data for 1000 users if they all have the same setting. The hash can be exchanged after the transformation has occured which is more easily atomic then manipulating lots of objects, potentially causing some data points to be updated and not others.
3. *Security* If you can pregenerate options you dont have to accept actual complex data objects from your users, less validation, less processor time spent validating aka harder to DDOS, more explicit contracts. 
4. *Processing* Generation can be done as a job at times when system load is low or in the background before a change. `Generate` could be fancified up quite a bit to understand priorities, restarts, parallel generation, inter key constraints, and probably other things.
5. *Jobs* It would be easy to add functionality to iterate over the data set and make necessary changes. Anything from disallowing a combination that is troublesome at the moment, to adding a new default value (which could be calculated off of other values), to cleaning up old settings, to anything else. Hashes can be exchanged as above in "Consistent datasets". It could be run as scheduled like "Processing" above.
6. *Defaults* Create some default objects and just spit out the default hash when theres no more appropriate answer, easy peasy.

Downsides? Splitting apart data, between finite and not really, which seems to represent one thing could be a bit confusing but its pretty hideable if you want to. Generating huge sets of data obviously takes a lot of time and adding just another key with a few values would pretty soon start meaning a hugely greater number of permutations to generate. Processors are fairly cheap however and you don't need to generate EVERYTHING if you don't want to. Additionally you should probably clean up old values. 

## How it works

There are five public functions: `Register`, `Generate`, `Store`, `Remove`, and `Get`.

`Register` allows you to build up the definition of an object that `Generate` can use later to pregenerate all permutations of finite data. You pass it the name of the object (so it can be kept track of internally), the key the data lives under and a list of all the value possibilities that exist for that key. You call Register per key that you need to set until you have created the object as you desire.

`Generate` creates every permutation of the `Register`ed data and `Store`s them.

`Store` takes an instance of an object, hashes it, writes it out to the file system and returns the hash.

`Remove` takes a hash and deletes it from the store.

`Get` takes a hash and returns the object from the file system if it exists.

Just using `Store` and `Get` will work but is just a partial implementation of a hashmap (not completely useless though). `Register` and `Generate` allow some interesting use cases like requiring a client not to send across the configuration object it's built up locally from user input but instead hash locally (need the exact same hash function on each side) and send the hash over. The server can "know" all the data that the hash represents by looking up the permutation. Less data has to be transferred, the hash can be validated / secured more easily, and still both sides can have an understanding of each other.

## Todos
* Don't just accept map[string]string, I haven't done go in a while and didn't want to mess around with how the types work
* Implement output so it would be pluggable - allow overriding output with a function that takes the hash and data and sticks it wherever
* Cleanup some of the code
* Rather then just passing an array of values you should be able to pass an iterator, type (like Boolean) and maybe other things?
* Allow creation of multiple instances that could store in different places / have different registered functions
