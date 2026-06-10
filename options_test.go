package dotpip

import (
	"testing"
)

func TestGenericOptions(t *testing.T) {
	var cmd CopyCommand
	WithDestination(nil)(&cmd)
	WithReplace()(&cmd)
	if cmd.Destination != nil || !cmd.Replace {
		t.Error("CopyCommand options failed")
	}

	var exp ExpireCommand
	WithExpireNX()(&exp)
	WithExpireXX()(&exp)
	WithExpireGT()(&exp)
	WithExpireLT()(&exp)
	if !exp.NX || !exp.XX || !exp.GT || !exp.LT {
		t.Error("ExpireCommand options failed")
	}

	var res RestoreCommand
	WithRestoreReplace()(&res)
	WithRestoreAbsTTL()(&res)
	WithRestoreIdleTime(10)(&res)
	WithRestoreFreq(5)(&res)
	if !res.Replace || !res.AbsTTL || res.IdleTime != 10 || res.Freq != 5 {
		t.Error("RestoreCommand options failed")
	}

	var scan ScanCommand
	WithScanMatch("*")(&scan)
	WithScanCount(10)(&scan)
	WithScanType(string(ObjectTypeString))(&scan)
	if scan.Match != "*" || scan.Count != 10 || scan.Type != string(ObjectTypeString) {
		t.Error("ScanCommand options failed")
	}

	var mig MigrateCommand
	WithMigrateCopy()(&mig)
	WithMigrateReplace()(&mig)
	WithMigrateKeys(Key{"a"})(&mig)
	if !mig.Copy || !mig.Replace || len(mig.Keys) != 1 {
		t.Error("MigrateCommand options failed")
	}
}

func TestGeoOptions(t *testing.T) {
	var add GeoAddCommand
	WithGeoAddNX()(&add)
	WithGeoAddXX()(&add)
	WithGeoAddCH()(&add)
	if !add.NX || !add.XX || !add.CH {
		t.Error("GeoAddCommand options failed")
	}

	var s GeoSearchCommand
	WithGeoSearchFromMember("m")(&s)
	WithGeoSearchFromLonLat(1.0, 2.0)(&s)
	WithGeoSearchByRadius(5.0, "km")(&s)
	WithGeoSearchByBox(10.0, 10.0, "m")(&s)
	WithGeoSearchAsc()(&s)
	WithGeoSearchDesc()(&s)
	WithGeoSearchCount(10, true)(&s)
	WithGeoSearchWithCoord()(&s)
	WithGeoSearchWithDist()(&s)
	WithGeoSearchWithHash()(&s)
	if s.FromMember != "m" || !s.UseLonLat || s.FromLongitude != 1.0 || s.FromLatitude != 2.0 || !s.UseRadius || s.ByRadius != 5.0 || s.RadiusUnit != "km" || !s.UseBox || s.ByBoxWidth != 10.0 || s.ByBoxHeight != 10.0 || s.BoxUnit != "m" || !s.Asc || !s.Desc || s.Count != 10 || !s.Any || !s.WithCoord || !s.WithDist || !s.WithHash {
		t.Error("GeoSearchCommand options failed")
	}

	var st GeoSearchStoreCommand
	WithGeoSearchStoreDist()(&st)
	if !st.StoreDist {
		t.Error("GeoSearchStoreCommand options failed")
	}
}

func TestHashesOptions(t *testing.T) {
	var h HRandFieldCommand
	WithHRandFieldWithValues()(&h)
	if !h.WithValues {
		t.Error("HRandFieldCommand options failed")
	}
}

func TestListsOptions(t *testing.T) {
	var l LPosCommand
	WithLPosRank(2)(&l)
	WithLPosCount(3)(&l)
	WithLPosMaxLen(10)(&l)
	if l.Rank != 2 || l.Count != 3 || l.MaxLen != 10 {
		t.Error("LPosCommand options failed")
	}
}

func TestStreamsOptions(t *testing.T) {
	var xa XAddCommand
	WithXAddNoMkStream()(&xa)
	WithXAddMaxLen(100, true)(&xa)
	WithXAddMinID("1-0", true)(&xa)
	WithXAddLimit(50)(&xa)
	if !xa.NoMkStream || xa.MaxLen != 100 || xa.MinID != "1-0" || !xa.Approx || xa.Limit != 50 {
		t.Error("XAddCommand options failed")
	}

	var xt XTrimCommand
	WithXTrimMaxLen(100, true)(&xt)
	WithXTrimMinID("1-0", true)(&xt)
	WithXTrimLimit(50)(&xt)
	if xt.MaxLen != 100 || !xt.Approx || xt.MinID != "1-0" || xt.Limit != 50 {
		t.Error("XTrimCommand options failed")
	}

	var xr XReadCommand
	WithXReadCount(10)(&xr)
	WithXReadBlock(1000)(&xr)
	if xr.Count != 10 || xr.Block != 1000 {
		t.Error("XReadCommand options failed")
	}

	var xrg XReadGroupCommand
	WithXReadGroupCount(10)(&xrg)
	WithXReadGroupBlock(1000)(&xrg)
	WithXReadGroupNoAck()(&xrg)
	if xrg.Count != 10 || xrg.Block != 1000 || !xrg.NoAck {
		t.Error("XReadGroupCommand options failed")
	}

	var xp XPendingCommand
	WithXPendingIdle(100)(&xp)
	WithXPendingRange("-", "+", 10)(&xp)
	if xp.Idle != 100 || xp.Start != "-" || xp.End != "+" || xp.Count != 10 {
		t.Error("XPendingCommand options failed")
	}

	var xc XClaimCommand
	WithXClaimIdle(100)(&xc)
	WithXClaimTime(1000)(&xc)
	WithXClaimRetryCount(3)(&xc)
	WithXClaimForce()(&xc)
	WithXClaimJustID()(&xc)
	if xc.Idle != 100 || xc.Time != 1000 || xc.RetryCount != 3 || !xc.Force || !xc.JustID {
		t.Error("XClaimCommand options failed")
	}

	var xac XAutoClaimCommand
	WithXAutoClaimCount(10)(&xac)
	WithXAutoClaimJustID()(&xac)
	if xac.Count != 10 || !xac.JustID {
		t.Error("XAutoClaimCommand options failed")
	}
}

func TestStringsOptions(t *testing.T) {
	var s SetCommand
	WithNX()(&s)
	WithXX()(&s)
	WithIfEq("a")(&s)
	WithIfNe("b")(&s)
	WithIfDeq("c")(&s)
	WithIfDne("d")(&s)
	WithGet()(&s)
	WithEx(10)(&s)
	WithPx(20)(&s)
	WithExAt(30)(&s)
	WithPxAt(40)(&s)
	WithKeepTTL()(&s)
	if !s.NX || !s.XX || s.IfEq != "a" || s.IfNe != "b" || s.IfDeq != "c" || s.IfDne != "d" || !s.Get || s.Ex != 10 || s.Px != 20 || s.ExAt != 30 || s.PxAt != 40 || !s.KeepTTL {
		t.Error("SetCommand options failed")
	}
}

func TestZSetsOptions(t *testing.T) {
	var za ZAddCommand
	WithZAddNX()(&za)
	WithZAddXX()(&za)
	WithZAddGT()(&za)
	WithZAddLT()(&za)
	WithZAddCH()(&za)
	WithZAddINCR()(&za)
	if !za.NX || !za.XX || !za.GT || !za.LT || !za.CH || !za.INCR {
		t.Error("ZAddCommand options failed")
	}

	var zr ZRangeCommand
	WithZRangeByScore()(&zr)
	WithZRangeByLex()(&zr)
	WithZRangeRev()(&zr)
	WithZRangeLimit(1, 2)(&zr)
	if !zr.ByScore || !zr.ByLex || !zr.Rev || zr.Offset != 1 || zr.Count != 2 {
		t.Error("ZRangeCommand options failed")
	}
}

func TestCommands(t *testing.T) {
	k1 := NewKey("a", "b", "c")
	if len(k1) != 3 || k1[0] != "a" || k1[1] != "b" || k1[2] != "c" {
		t.Error("NewKey with multiple args failed")
	}

	k2 := NewKey("a:b:c")
	if len(k2) != 3 || k2[0] != "a" || k2[1] != "b" || k2[2] != "c" {
		t.Error("NewKey with delimiter failed")
	}

	k3 := NewKeyWithDelimiter("a:b:c", ":")
	if len(k3) != 3 || k3[0] != "a" || k3[1] != "b" || k3[2] != "c" {
		t.Error("NewKeyWithDelimiter string, delimiter failed")
	}

	k4 := NewKeyWithDelimiter("a", "b", "c", ".")
	if len(k4) != 3 || k4[0] != "a" || k4[1] != "b" || k4[2] != "c" {
		t.Error("NewKeyWithDelimiter multiple args, delimiter failed")
	}

	k5 := NewKeyWithDelimiter("a.b", "c", ".")
	if len(k5) != 3 || k5[0] != "a" || k5[1] != "b" || k5[2] != "c" {
		t.Error("NewKeyWithDelimiter mixed args, delimiter failed")
	}

	dp := New(nil)
	if dp != nil {
		t.Error("New failed")
	}

	f := &DataTypeFormatter{
		JSONEncode: func(value any) (any, error) {
			return value, nil
		},
	}
	res, err := f.JSONSetEncode("test")
	if err != nil || res != "test" {
		t.Error("JSONSetEncode failed")
	}
}

func TestScriptFlushOptions(t *testing.T) {
	var cmd ScriptFlushCommand
	WithScriptFlushSync()(&cmd)
	WithScriptFlushAsync()(&cmd)

	if !cmd.Sync || !cmd.Async {
		t.Errorf("Expected Sync and Async to be true")
	}
}

func TestVectorOptionsCoverage(t *testing.T) {
	WithVAddCas(true)(&VAddOptions{})
	WithVAddEF(1)(&VAddOptions{})
	WithVAddM(1)(&VAddOptions{})
	WithVAddQuant("q")(&VAddOptions{})
	WithVAddReduceDim(1)(&VAddOptions{})

	WithVSimCount(1)(&VSimOptions{})
	WithVSimEF(1)(&VSimOptions{})
	WithVSimEpsilon(1.0)(&VSimOptions{})
	WithVSimFilter("test")(&VSimOptions{})
	WithVSimFilterEF(1)(&VSimOptions{})
	WithVSimNoThread(true)(&VSimOptions{})
	WithVSimTruth(true)(&VSimOptions{})
	WithVSimWithAttribs(true)(&VSimOptions{})
	WithVSimWithScores(true)(&VSimOptions{})
}
