//Package response contains model classes wrapping the response receved by a ProbeClient
package response
//ProbeResponse is a model class for response of a probe request
type ProbeResponse interface{
    //GetType returns the ProbeType
    //Ideally this value should be the same as the ProbeClient's ProbeType value
    GetType()string
    //AsMap converts the response into a Map.
    //This map could be marshalled & sent (to RuleEngine) for computing metrics
    AsMap()map[string]interface{}
}